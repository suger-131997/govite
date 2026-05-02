// Package govite provides integration between Go HTTP servers and Vite,
// enabling server-side rendering with type-safe props passing from Go to
// frontend components, and automatic TypeScript type definition generation.
//
// # Overview
//
// govite bridges a Go HTTP backend with a Vite-powered frontend. It supports
// two rendering modes:
//   - Development mode: proxies asset requests to the Vite dev server and
//     generates entry point files on-the-fly.
//   - Production mode: uses the Vite build manifest to inject the correct
//     hashed asset URLs into the HTML response.
//
// # Usage
//
// Use [WithRenderCreatorForDev] or [WithRenderCreatorForProd] to attach a
// renderer factory to a [context.Context], then create page handlers with
// [NewPageHandler] and obtain their [net/http.HandlerFunc] via
// [PageHandler.Handler].
//
// Page-specific values such as the HTTP status code and title can be passed
// through the context using [WithStatusCode] and [WithTitle].
package govite
