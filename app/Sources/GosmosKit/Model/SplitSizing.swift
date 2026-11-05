import Foundation

/// A split's first pane, including windows smaller than the panes' preferred minimums.
public enum SplitSizing {
    public static func firstLength(total: Double, fraction: Double, minFirst: Double, minSecond: Double) -> Double {
        guard total.isFinite, total > 1 else { return 0 }
        let available = total - 1
        let first = minFirst.isFinite ? max(0, minFirst) : 0
        let second = minSecond.isFinite ? max(0, minSecond) : 0
        let minimum = first + second
        // Both panes shrink together instead of pushing the second out of the window.
        if minimum > available { return available * (first / minimum) }
        let share = fraction.isFinite ? min(1, max(0, fraction)) : 0.5
        return min(max(available * share, first), available - second)
    }
}
