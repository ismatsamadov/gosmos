import Foundation

/// A lesson in a library. Lesson uuids are unique within their library, so the learner
/// database keys everything of yours by both (library_uuid, lesson_uuid).
public struct LessonRef: Hashable, Sendable {
    public let library: String
    public let lesson: String

    public init(library: String, lesson: String) {
        self.library = library
        self.lesson = lesson
    }
}

/// A library's identity, from its single library_info row.
public struct LibraryInfo: Hashable, Sendable {
    public let uuid: String
    public let title: String

    public init(uuid: String, title: String) {
        self.uuid = uuid
        self.title = title
    }
}

/// A library you study in, as the learner database lists it (its library table).
public struct LearnerLibrary: Hashable, Sendable {
    public let uuid: String
    public let title: String
    /// `builtin:<file name>` for a library built into the app, else the file's path.
    public let location: String
    public let addedAt: String

    public init(uuid: String, title: String, location: String, addedAt: String) {
        self.uuid = uuid
        self.title = title
        self.location = location
        self.addedAt = addedAt
    }

    public var isBuiltIn: Bool { location.hasPrefix(BuiltInLibraries.prefix) }
}

/// The course libraries built into the app: Resources/Courses/*.sqlite.
public enum BuiltInLibraries {
    static let prefix = "builtin:"

    /// Where the app keeps them.
    public static var folder: URL? { Bundle.main.resourceURL?.appendingPathComponent("Courses", isDirectory: true) }

    /// The libraries in `folder`, by file name; none when there is no folder.
    public static func files(in folder: URL?) -> [URL] {
        guard let folder, let names = try? FileManager.default.contentsOfDirectory(atPath: folder.path) else { return [] }
        return names.filter { $0.hasSuffix(".sqlite") && !$0.hasPrefix(".") }.sorted()
            .map { folder.appendingPathComponent($0, isDirectory: false) }
    }

    /// The learner database's location for a built-in library.
    public static func location(of file: URL) -> String { prefix + file.lastPathComponent }

    /// The libraries you opened yourself, in the order you added them: the learner
    /// database's libraries that are not built in. A built-in one comes from the app instead,
    /// and one no longer shipped stays out of sight with its progress kept.
    public static func opened(from libraries: [LearnerLibrary]) -> [LearnerLibrary] {
        libraries.filter { !$0.isBuiltIn }.sorted { ($0.addedAt, $0.uuid) < ($1.addedAt, $1.uuid) }
    }
}
