// Draws Gosmos's icon at every size of a macOS iconset: a Go Blue planet inside a tilted
// orbit, with a small satellite on it, on a night-blue tile in the macOS rounded square.
// Layers: the tile and its shadow, the sky with a few stars, the planet's glow, the back of
// the orbit, the planet, the front of the orbit, the satellite and a light rim. Each size is
// drawn on its own pixel grid, with fewer details at 16 and 32 px, so it stays sharp.
// Usage: make-icon <iconset folder>   (scripts/make-icon.sh runs it)

import AppKit
import SwiftUI

func rgb(_ hex: UInt32, _ alpha: CGFloat = 1) -> CGColor {
    CGColor(srgbRed: CGFloat(hex >> 16 & 0xFF) / 255, green: CGFloat(hex >> 8 & 0xFF) / 255,
            blue: CGFloat(hex & 0xFF) / 255, alpha: alpha)
}

let goBlue = rgb(0x00ADD8)
let goBlueLight = rgb(0x8FE3F5)
let goBlueDeep = rgb(0x005A78)
let nightTop = rgb(0x1B2E62)
let nightBottom = rgb(0x0A1230)
let ice = rgb(0xEAF8FD)
let sRGB = CGColorSpace(name: CGColorSpace.sRGB)!

/// The macOS icon grid's rounded square, with Apple's continuous corners.
func squircle(_ rect: CGRect) -> CGPath {
    RoundedRectangle(cornerRadius: rect.width * 0.225, style: .continuous).path(in: rect).cgPath
}

/// Stars of the sky, in the 1024 grid: x, y, radius, opacity. They keep clear of the planet.
let stars: [(CGFloat, CGFloat, CGFloat, CGFloat)] = [
    (238, 792, 7, 0.95), (336, 702, 4, 0.55), (452, 842, 5, 0.7), (622, 868, 4, 0.5),
    (812, 694, 4, 0.55), (184, 560, 4, 0.45), (868, 448, 5, 0.6), (232, 268, 4, 0.5),
    (766, 214, 5, 0.6), (420, 176, 3, 0.45), (884, 812, 3, 0.45),
]

func render(pixels n: Int) -> NSBitmapImageRep {
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: n, pixelsHigh: n, bitsPerSample: 8, samplesPerPixel: 4,
                               hasAlpha: true, isPlanar: false, colorSpaceName: .deviceRGB, bytesPerRow: 0, bitsPerPixel: 0)!
    rep.size = NSSize(width: n, height: n)
    NSGraphicsContext.saveGraphicsState()
    NSGraphicsContext.current = NSGraphicsContext(bitmapImageRep: rep)
    let ctx = NSGraphicsContext.current!.cgContext
    ctx.setShouldAntialias(true)
    let s = CGFloat(n) / 1024
    let small = n <= 32
    let tiny = n <= 16
    // Snap to whole pixels at small sizes, so edges are crisp.
    func p(_ v: CGFloat) -> CGFloat { small ? (v * s).rounded() : v * s }
    func circle(_ x: CGFloat, _ y: CGFloat, _ r: CGFloat) -> CGRect {
        CGRect(x: p(x - r), y: p(y - r), width: p(x + r) - p(x - r), height: p(y + r) - p(y - r))
    }

    // The tile, with a soft shadow under it as macOS icons have.
    let tile = CGRect(x: p(100), y: p(100), width: p(924) - p(100), height: p(924) - p(100))
    let shape = squircle(tile)
    if n >= 64 {
        ctx.saveGState()
        ctx.setShadow(offset: CGSize(width: 0, height: -10 * s), blur: 24 * s, color: NSColor.black.withAlphaComponent(0.3).cgColor)
        ctx.addPath(shape)
        ctx.setFillColor(nightBottom)
        ctx.fillPath()
        ctx.restoreGState()
    }
    ctx.saveGState()
    ctx.addPath(shape)
    ctx.clip()

    // Layer: the night sky, lighter at the top.
    let sky = CGGradient(colorsSpace: sRGB, colors: [nightTop, nightBottom] as CFArray, locations: [0, 1])!
    ctx.drawLinearGradient(sky, start: CGPoint(x: 0, y: tile.maxY), end: CGPoint(x: 0, y: tile.minY), options: [])

    // Layer: stars, only where they can be more than a smudge.
    if n >= 128 {
        for (x, y, r, a) in stars {
            ctx.setFillColor(rgb(0xFFFFFF, a))
            ctx.fillEllipse(in: circle(x, y, max(r, 1.2 / s)))
        }
    }

    // The planet and its orbit, in the 1024 grid; bolder and level at small sizes.
    let center = CGPoint(x: 512, y: small ? 512 : 496)
    let radius: CGFloat = tiny ? 262 : small ? 250 : 214
    let tilt: CGFloat = small ? 0 : -0.34
    let rx: CGFloat = tiny ? 372 : small ? 380 : 372
    let ry: CGFloat = tiny ? 118 : small ? 112 : 104
    let ringWidth: CGFloat = tiny ? 1.4 : small ? 2.2 : max(30 * s, 2.4) // in pixels
    var orbit = CGAffineTransform(translationX: center.x * s, y: center.y * s).rotated(by: tilt)
    let ring = CGPath(ellipseIn: CGRect(x: -rx * s, y: -ry * s, width: 2 * rx * s, height: 2 * ry * s), transform: &orbit)

    // Layer: the planet's glow on the sky.
    if n >= 64 {
        let glow = CGGradient(colorsSpace: sRGB, colors: [rgb(0x00ADD8, 0.32), rgb(0x00ADD8, 0)] as CFArray, locations: [0, 1])!
        let c = CGPoint(x: center.x * s, y: center.y * s)
        ctx.drawRadialGradient(glow, startCenter: c, startRadius: radius * s * 0.9, endCenter: c, endRadius: radius * s * 1.85, options: [])
    }

    // Layer: the back of the orbit, behind the planet and a little dimmer.
    ctx.addPath(ring)
    ctx.setLineWidth(ringWidth)
    ctx.setStrokeColor(rgb(0xEAF8FD, small ? 0.75 : 0.5))
    ctx.strokePath()

    // Layer: the planet, lit from the upper left.
    let planet = circle(center.x, center.y, radius)
    ctx.saveGState()
    ctx.addEllipse(in: planet)
    ctx.clip()
    if tiny {
        ctx.setFillColor(goBlue)
        ctx.fill(planet)
    } else {
        let light = CGPoint(x: planet.midX - planet.width * 0.2, y: planet.midY + planet.height * 0.24)
        let shading = CGGradient(colorsSpace: sRGB, colors: [goBlueLight, goBlue, goBlueDeep] as CFArray, locations: [0, 0.5, 1])!
        ctx.drawRadialGradient(shading, startCenter: light, startRadius: 0, endCenter: light, endRadius: planet.width * 0.78,
                               options: [.drawsAfterEndLocation])
        if n >= 128 {
            // Two faint bands across the planet, following the orbit's tilt.
            ctx.concatenate(orbit)
            ctx.setFillColor(rgb(0xFFFFFF, 0.1))
            ctx.fill(CGRect(x: -radius * s, y: 46 * s, width: 2 * radius * s, height: 30 * s))
            ctx.fill(CGRect(x: -radius * s, y: -96 * s, width: 2 * radius * s, height: 22 * s))
        }
    }
    ctx.restoreGState()

    // Layer: the front of the orbit, across the planet: the ring again, on the near half only.
    ctx.saveGState()
    ctx.concatenate(orbit)
    ctx.clip(to: CGRect(x: -(rx + 40) * s, y: -(ry + 40) * s, width: 2 * (rx + 40) * s, height: (ry + 40) * s))
    ctx.concatenate(orbit.inverted())
    ctx.addPath(ring)
    ctx.setLineWidth(ringWidth)
    ctx.setStrokeColor(ice)
    ctx.strokePath()
    ctx.restoreGState()

    // Layer: the satellite, a small bright body on the near side of the orbit; too small to
    // tell from the ring at 16 and 32 px.
    if !small {
        let angle: CGFloat = -0.16 * .pi
        let at = CGPoint(x: rx * cos(angle) * s, y: ry * sin(angle) * s).applying(orbit)
        let r = max(30 * s, 2.2)
        let halo = CGGradient(colorsSpace: sRGB, colors: [rgb(0x8FE3F5, 0.7), rgb(0x8FE3F5, 0)] as CFArray, locations: [0, 1])!
        ctx.drawRadialGradient(halo, startCenter: at, startRadius: r * 0.8, endCenter: at, endRadius: r * 2.6, options: [])
        ctx.setFillColor(rgb(0xFFFFFF))
        ctx.fillEllipse(in: CGRect(x: at.x - r, y: at.y - r, width: 2 * r, height: 2 * r))
    }

    // A thin light rim along the edge, for a glassy edge.
    if n >= 128 {
        ctx.addPath(shape)
        ctx.setStrokeColor(rgb(0xFFFFFF, 0.12))
        ctx.setLineWidth(6 * s)
        ctx.strokePath()
    }
    ctx.restoreGState()
    NSGraphicsContext.restoreGraphicsState()
    return rep
}

let args = CommandLine.arguments
guard args.count == 2 else {
    FileHandle.standardError.write(Data("usage: make-icon <iconset folder>\n".utf8))
    exit(2)
}
let folder = URL(fileURLWithPath: args[1], isDirectory: true)
try FileManager.default.createDirectory(at: folder, withIntermediateDirectories: true)
for base in [16, 32, 128, 256, 512] {
    for scale in [1, 2] {
        let pixels = base * scale
        let name = scale == 1 ? "icon_\(base)x\(base).png" : "icon_\(base)x\(base)@2x.png"
        guard let png = render(pixels: pixels).representation(using: .png, properties: [:]) else { exit(1) }
        try png.write(to: folder.appendingPathComponent(name))
    }
}
