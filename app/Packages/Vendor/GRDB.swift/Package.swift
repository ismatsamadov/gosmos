// swift-tools-version:6.1
// GRDB.swift 7.11.1 (b83108d), kept in Gosmos's repository so the app builds offline. This is
// the package's own manifest with only what Gosmos builds: the GRDB library on the system's
// SQLite. Its tests, documentation and the SQLCipher and custom-SQLite variants are left out.

import PackageDescription

let swiftSettings: [SwiftSetting] = [
    .define("SQLITE_ENABLE_FTS5"),
    .define("SQLITE_ENABLE_SNAPSHOT"),
    // Not all Linux distributions have support for WAL snapshots.
    .define("SQLITE_DISABLE_SNAPSHOT", .when(platforms: [.linux])),
]

let package = Package(
    name: "GRDB",
    platforms: [
        .iOS(.v13),
        .macOS(.v10_15),
        .tvOS(.v13),
        .watchOS(.v7),
    ],
    products: [
        .library(name: "GRDBSQLite", targets: ["GRDBSQLite"]),
        .library(name: "GRDB", targets: ["GRDB"]),
    ],
    targets: [
        .systemLibrary(
            name: "GRDBSQLite",
            providers: [.apt(["libsqlite3-dev"])]),
        .target(
            name: "GRDB",
            dependencies: [
                .target(name: "GRDBSQLite"),
            ],
            path: "GRDB",
            resources: [.copy("PrivacyInfo.xcprivacy")],
            swiftSettings: swiftSettings + [
                .enableUpcomingFeature("MemberImportVisibility"),
                // Its warnings are not Gosmos's to fix; they stay quiet, as when SwiftPM fetched it.
                .unsafeFlags(["-suppress-warnings"]),
            ]),
    ],
    swiftLanguageModes: [.v6]
)
