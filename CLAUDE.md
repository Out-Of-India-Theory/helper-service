# CLAUDE.md

A single-purpose utility service (Go + Gin, base path `helper-service`, port `:8080`): it **renders localised push-notification images for supply (astrologer) assignment** using a headless Chromium. It has no database. Still on the `initial-commit` branch — treat it as young code, not a mature service.

## The only real endpoint
`POST helper-service/pn-image/:supply_id` → `controller/image_generator` → `service/image_generator.GenerateImage(ctx, supplyId)`.

Pipeline per request:
1. `service/supply` GETs `{OMSClientConfig.Address}/internal/supply/supplies/{supply_id}` for the supply's `NameV1` (per-language name map), `Languages`, `ImageWithoutBackground` and experience.
2. Downloads the supply's transparent-background photo and base64-encodes it.
3. **For each language the supply serves**, fills `assets/template.html` (1280×800) with the language's font path, `assets/images/background.png`, the base64 photo, and the translated copy, then screenshots it with **chromedp** (headless Chromium).
4. `service/image_uploader` POSTs the PNG bytes to platform-service's document upload endpoint and returns the URL.

So one call produces **N images, one per language** — a supply serving 5 languages costs 5 Chromium screenshots.

## Things that will bite you
- **Translations are hardcoded Go string maps**, not the commons `language` repo. `service/image_generator/image_generator_service.go` holds `prashnaTranslations`, `jyotishaTranslations`, `experienceText`, `checkStatusText` for `en, hi, kn, gu, ta, te, mr`. Adding a language means editing four maps **and** `fontMap` **and** shipping the font. A language present in the supply's `Languages` but missing from `fontMap` yields an empty font path and a broken render.
- `prashnaTranslations` and `jyotishaTranslations` are currently **identical** — if the two flows are supposed to differ, that's an open bug.
- **Fonts live in two places.** `assets/fonts/*.ttf` is committed and referenced by absolute path at runtime (`filepath.Abs`), while the Dockerfile *also* curls a different, partly non-overlapping set into `/usr/share/fonts/noto` (it fetches Malayalam and Bengali, which the code never uses, and Telugu-**Light** rather than Regular). The code path uses the committed `assets/` files, so **`filepath.Abs` means the process must run with the repo root as its working directory.**
- **The upload target is wrong-looking but load-bearing:** `image_uploader` builds `{OMSClientConfig.Address}/platform/document/v1/upload/...` — it uses the **OMS** config address for a **platform-service** path. That works only because both point at the same internal load balancer. `PlatformClientConfig` exists in config but is unused. Don't "fix" the URL without checking the deployed addresses.
- The upload use-case name differs by environment: `prod_supply_pn_images` when `ServerConfig.Env` is `PROD`/`PRODUCTION`, else `jyotisha_pn_image`. Both must exist in platform-service's `DocumentStoreConfig`.
- **One shared Chromium context** is created once in `InitImageGeneratorService` (`headless`, `no-sandbox`, `disable-dev-shm-usage`) and reused for every request — it is not per-request isolated, so concurrent renders share browser state and a crashed browser stays crashed until restart.
- Both HTTP clients are hand-rolled `net/http` (not commons `http_client`), and `service/supply` uses a **10-minute** timeout while its `OMSClientConfig.Timeout` is `10s` — the config value is ignored.
- `config.json` in the repo points at the **staging internal LB** for both client blocks.

## Build / run
```bash
go build ./... && go vet ./...
go run .        # from repo root — asset paths are resolved relative to CWD
```
Needs a local Chromium/Chrome for `chromedp`. No tests. Go 1.24.3, commons `v0.1.0-rc`.

The Dockerfile is unusual for this fleet: it is a **single-stage `golang:1.24-alpine` image** (no distroless/scratch stage), installs `chromium` and symlinks it to `/usr/bin/google-chrome`, runs `go mod tidy` at build time (so a build can pick up dependency drift), and cross-builds `GOOS=linux GOARCH=arm64`. The final image ships the whole Go toolchain — large, and worth knowing before optimising.

Deploy: Docker → ECR → ECS (`ap-south-1`) via `.github/workflows/aws.yml`.

## Files
- `service/image_generator/image_generator_service.go` — everything: translations, fonts, HTML fill, chromedp screenshot.
- `service/supply/` — OMS supply fetch + `models.go` response shape.
- `service/image_uploader/` — S3 upload via platform-service.
- `assets/template.html` — the render template; placeholders are `{{FONT_PATH}}`, `{{BG_PATH}}` and friends, substituted by string replacement (not `html/template`).

## If you extend this service
The name suggests a grab-bag, and the README says "first use case". Before adding an unrelated helper here, consider whether it belongs in the owning service instead — this repo has no DB, no auth, and a Chromium dependency that makes it heavy to scale. Follow the existing shape (`controller/<x>` → `service/facade` → `service/<x>`) and register the route in `server/router.go`. Conventions match the fleet: `context.Context` first, `fmt.Errorf("Func: %w", err)`, zap via `logging.WithContext(ctx)`.
