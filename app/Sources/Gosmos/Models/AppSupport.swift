import Foundation
import GosmosKit

/// A failure the app shows: as an inline banner, an alert, or as the window's content when
/// your study data can't be opened.
enum AppError: LocalizedError, Identifiable, Equatable {
    case database(AppDatabaseError)
    case learner(LearnerDatabaseError)
    case lesson(String)
    case files(String)

    var id: String { errorDescription ?? "error" }

    var errorDescription: String? {
        switch self {
        case .database(let e): e.errorDescription
        case .learner(let e): e.errorDescription
        case .lesson(let m): m
        case .files(let m): m
        }
    }

    var recoverySuggestion: String? {
        switch self {
        case .database(let e): e.recoverySuggestion
        case .learner(let e): e.recoverySuggestion
        case .lesson, .files: nil
        }
    }

    /// The headline for the error.
    var title: String {
        switch self {
        case .database(let e): e.title
        case .learner(let e): e.title
        case .lesson: "The Lesson Cannot Open"
        case .files: "Files Error"
        }
    }
}

/// Settings that belong to this Mac, in UserDefaults; your progress lives in the learner
/// database.
struct LocalPreferences {
    static let goPathKey = "GoPath"
    static let unsettledKey = "UnsettledCode"
    /// Scripted runs (screenshots, tests of the app) keep their study data apart:
    /// `-GosmosDataFolder <folder>` holds the learner database and the lesson folders. It lives
    /// in the argument domain only, never stored.
    static let dataFolderKey = "GosmosDataFolder"

    let defaults: UserDefaults

    var goPath: String? {
        let p = defaults.string(forKey: LocalPreferences.goPathKey)?.trimmingCharacters(in: .whitespaces)
        return (p?.isEmpty ?? true) ? nil : p
    }

    func setGoPath(_ path: String?) { defaults.set(path ?? "", forKey: LocalPreferences.goPathKey) }

    var dataFolder: URL? {
        defaults.string(forKey: LocalPreferences.dataFolderKey).map { URL(fileURLWithPath: ($0 as NSString).expandingTildeInPath, isDirectory: true) }
    }

    /// Code files left in conflict when the app quit, as "<library>/<lesson>/<file>": the
    /// lesson folder on this Mac holds the other editor's version, so this list belongs to the
    /// Mac too.
    func unsettledFiles(in ref: LessonRef) -> Set<String> {
        let prefix = ref.library + "/" + ref.lesson + "/"
        return Set((defaults.stringArray(forKey: LocalPreferences.unsettledKey) ?? []).compactMap { entry in
            entry.hasPrefix(prefix) ? String(entry.dropFirst(prefix.count)) : nil
        })
    }

    func setUnsettled(_ file: String, in ref: LessonRef, _ unsettled: Bool) {
        let entry = ref.library + "/" + ref.lesson + "/" + file
        var list = (defaults.stringArray(forKey: LocalPreferences.unsettledKey) ?? []).filter { $0 != entry }
        if unsettled { list.append(entry) }
        if list.isEmpty { defaults.removeObject(forKey: LocalPreferences.unsettledKey) } else { defaults.set(list, forKey: LocalPreferences.unsettledKey) }
    }
}

/// The app's read-only resources. Flow tests supply their own tiny libraries and schema.
struct AppResources {
    let libraries: URL?
    let learnerSchema: URL?
    let reference: URL?

    static var bundled: Self {
        Self(libraries: BuiltInLibraries.folder, learnerSchema: LearnerDatabase.bundledSchema,
             reference: ReferenceDatabase.bundledURL)
    }
}
