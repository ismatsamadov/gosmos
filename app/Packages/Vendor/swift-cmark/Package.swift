// swift-tools-version:5.3
// swift-cmark 0.9.0 (08ddb52), kept in Gosmos's repository so the app builds offline. This is
// the package's own manifest with only the two libraries swift-markdown uses; the command-line
// program, its API tests and the build files for other systems are left out.

import PackageDescription

#if os(Windows)
    // When building the library on Windows, we do not have proper control of
    // whether it is being built statically or dynamically.  However, given the
    // current default of static, we will assume that we are building
    // statically.  More importantly, should this not hold, this will fail at
    // link time.
    let cSettings: [CSetting] = [
        .define("CMARK_GFM_STATIC_DEFINE", .when(platforms: [.windows])),
    ]
#else
    let cSettings: [CSetting] = []
#endif

let package = Package(
    name: "cmark-gfm",
    products: [
        .library(
            name: "cmark-gfm",
            targets: ["cmark-gfm"]),
        .library(
            name: "cmark-gfm-extensions",
            targets: ["cmark-gfm-extensions"]),
    ],
    targets: [
        .target(name: "cmark-gfm",
          path: "src",
          exclude: [
            "scanners.re",
            "libcmark-gfm.pc.in",
            "config.h.in",
            "CMakeLists.txt",
          ],
          cSettings: cSettings
        ),
        .target(name: "cmark-gfm-extensions",
          dependencies: [
            "cmark-gfm",
          ],
          path: "extensions",
          exclude: [
            "CMakeLists.txt",
            "ext_scanners.re",
          ],
          cSettings: cSettings
        ),
    ]
)
