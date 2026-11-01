import Foundation

/// Two versions of a file side by side, line by line, as Show Solution sets your code beside
/// the lesson's solution: lines that match sit together; between them, the lines that differ
/// are paired in order, and the rest stand alone. Spacing alone is no difference.
public enum LineDiff {
    /// One row: a line of the first version, of the second, or of both, by index.
    public struct Row: Hashable, Sendable {
        public let left: Int?
        public let right: Int?
        public let differs: Bool

        public init(left: Int?, right: Int?, differs: Bool) {
            self.left = left
            self.right = right
            self.differs = differs
        }
    }

    /// The rows for two versions' lines. Long files take a while, so it is worked out off the
    /// main actor, and it gives up once its task is cancelled (the sheet showing it closed).
    public static func rows(_ left: [String], _ right: [String]) throws(CancellationError) -> [Row] {
        var rows: [Row] = []
        var removed: [Int] = [], added: [Int] = []
        func pairUp() {
            for k in 0..<max(removed.count, added.count) {
                rows.append(Row(left: k < removed.count ? removed[k] : nil, right: k < added.count ? added[k] : nil, differs: true))
            }
            removed = []
            added = []
        }
        for step in try steps(left.map(comparable), right.map(comparable)) {
            switch step {
            case .same(let l, let r):
                pairUp()
                rows.append(Row(left: l, right: r, differs: false))
            case .removed(let l): removed.append(l)
            case .added(let r): added.append(r)
            }
        }
        pairUp()
        return rows
    }

    /// A line as compared: runs of spaces as one, none at the ends.
    static func comparable(_ line: String) -> String {
        line.split(whereSeparator: \.isWhitespace).joined(separator: " ")
    }

    enum Step: Equatable {
        case same(Int, Int)
        case removed(Int)
        case added(Int)
    }

    /// The longest common run of lines, as steps through both. Past four million line pairs
    /// (two files of 2,000 lines), it compares line by line instead. Cancelled, it stops
    /// within a row of the table.
    static func steps(_ a: [String], _ b: [String]) throws(CancellationError) -> [Step] {
        let n = a.count, m = b.count
        guard n * m <= 4_000_000 else {
            return (0..<max(n, m)).flatMap { i -> [Step] in
                if i < n, i < m { return a[i] == b[i] ? [.same(i, i)] : [.removed(i), .added(i)] }
                return i < n ? [.removed(i)] : [.added(i)]
            }
        }
        // lcs[i][j]: the longest common subsequence of a[i...] and b[j...].
        var lcs = [[Int32]](repeating: [Int32](repeating: 0, count: m + 1), count: n + 1)
        if n > 0 && m > 0 {
            for i in stride(from: n - 1, through: 0, by: -1) {
                if Task.isCancelled { throw CancellationError() }
                for j in stride(from: m - 1, through: 0, by: -1) {
                    lcs[i][j] = a[i] == b[j] ? lcs[i + 1][j + 1] + 1 : max(lcs[i + 1][j], lcs[i][j + 1])
                }
            }
        }
        var steps: [Step] = []
        var i = 0, j = 0
        while i < n && j < m {
            if a[i] == b[j] {
                steps.append(.same(i, j))
                i += 1
                j += 1
            } else if lcs[i + 1][j] >= lcs[i][j + 1] {
                steps.append(.removed(i))
                i += 1
            } else {
                steps.append(.added(j))
                j += 1
            }
        }
        while i < n { steps.append(.removed(i)); i += 1 }
        while j < m { steps.append(.added(j)); j += 1 }
        return steps
    }
}
