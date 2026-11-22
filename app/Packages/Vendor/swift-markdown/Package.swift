// swift-tools-version:6.2
/*
 This source file is part of the Swift.org open source project

 Copyright (c) 2021-2023 Apple Inc. and the Swift project authors
 Licensed under Apache License v2.0 with Runtime Library Exception

 See https://swift.org/LICENSE.txt for license information
 See https://swift.org/CONTRIBUTORS.txt for Swift project authors
*/

// swift-markdown 0.9.0 (25cb61d), kept in Gosmos's repository so the app builds offline. This
// is the package's own manifest with only the Markdown library, on the swift-cmark beside it
// rather than fetched; its tests and the DocC plugin are left out.

import PackageDescription

let package = Package(
    name: "swift-markdown",
    products: [
        .library(
            name: "Markdown",
            targets: ["Markdown"]),
    ],
    dependencies: [
        .package(path: "../swift-cmark"),
    ],
    targets: [
        .target(
            name: "Markdown",
            dependencies: [
                "CAtomic",
                .product(name: "cmark-gfm", package: "swift-cmark"),
                .product(name: "cmark-gfm-extensions", package: "swift-cmark"),
            ],
            exclude: [
                "CMakeLists.txt"
            ],
            swiftSettings: [
                .unsafeFlags(["-Xcc", "-DCMARK_GFM_STATIC_DEFINE"], .when(platforms: [.windows])),
                // Its warnings are not Gosmos's to fix; they stay quiet, as when SwiftPM fetched it.
                .unsafeFlags(["-suppress-warnings"]),
            ]
        ),
        .target(name: "CAtomic"),
    ],
    swiftLanguageModes: [.v5]
)
