<img src="./assets/gpxl_logo.png" alt="GPxl logo" width="300" style="display: block; margin: auto;" />

# GPxl

GPxl is a Go-based image processing package built to explore **pixel-level transformations**, **image algorithms**, and **Go’s standard image APIs** through first-principles implementations.

Rather than wrapping existing image-processing libraries, GPxl focuses on implementing filters and geometric transformations directly to better understand how image pipelines work under the hood.

---

## Overview

GPxl is a learning-focused image processing library written in Go.  
The project prioritizes **clarity, correctness, and understanding** over production-level optimization.

Key areas of exploration include:
- Pixel-wise image manipulation
- Image representations in Go (`image.Image`, `RGBA`, etc.)
- Algorithmic tradeoffs in CPU-bound workloads
- API and package design in Go

---

## Learning Goals

- Understand how images are represented and manipulated in Go
- Implement common image filters from first principles
- Explore geometric transformations at the pixel level
- Practice composable API design
- Reason about performance and concurrency tradeoffs

---

## Features

### Transformations
- [x] Reflect (Vertical & Horizontal)
- [x] Rotate (partial implementation)
- [ ] Resize
- [ ] Crop
- [ ] Shear

### Filters

#### Color
- [x] Grayscale
- [x] Sepia
- [x] Duotone
- [x] Cool
- [x] Warm
- [ ] Invert

#### Blends
- [ ] Darken
- [ ] Multiply

#### Enhance
- [ ] Saturation
- [ ] White Balance
- [ ] Curves
- [ ] Levels

### Effects
- [ ] Blur
- [ ] Edge Detection
- [ ] Emboss
- [ ] Sharpen
- [ ] Sobel
- [ ] Threshold

---

## Concurrency (Learning Focus)

GPxl explores Go concurrency by parallelizing per-pixel operations.

Approach:
- Compute `n = runtime.NumCPU()`
- Split the image into `n` row chunks
- Process each chunk in a goroutine
- Wait for completion via `sync.WaitGroup`
- Each goroutine writes to a distinct row range in the destination image to avoid overlapping writes.

This was built to learn:
- goroutines + WaitGroups
- safe work partitioning (no shared writes across chunks)
- practical limits of parallelism on CPU-bound loops

---

## Tradeoffs & Non-Goals

- GPxl does not aim to compete with optimized native image libraries
- SIMD, GPU acceleration, and CGO are intentionally avoided
- Performance optimizations are secondary to readability and correctness
- The project focuses on single-image transformations rather than batch pipelines

---

## Project Status

GPxl is an active learning project.

Planned areas of exploration:
- Additional convolution-based filters
- Improved transformation composition
- Benchmarking and profiling
- Further concurrency experimentation

---

## What I’ve Learned So Far

- How Go represents images and pixel buffers internally
- The cost of per-pixel operations in CPU-bound code
- How design decisions affect composability and testability
- Where concurrency helps — and where it doesn’t

