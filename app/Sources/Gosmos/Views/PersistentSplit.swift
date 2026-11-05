import GosmosKit
import SwiftUI

/// Two panes with a draggable divider whose position is remembered in UserDefaults (a window
/// choice that belongs to this Mac). SwiftUI's split views do not keep their positions.
struct PersistentSplit<First: View, Second: View>: View {
    let axis: Axis
    let minFirst: CGFloat
    let minSecond: CGFloat
    @AppStorage private var fraction: Double
    @ViewState private var dragStart: Double?
    private let first: First
    private let second: Second

    init(_ axis: Axis, key: String, defaultFraction: Double, minFirst: CGFloat, minSecond: CGFloat,
         @ViewBuilder first: () -> First, @ViewBuilder second: () -> Second) {
        self.axis = axis
        self.minFirst = minFirst
        self.minSecond = minSecond
        _fraction = AppStorage(wrappedValue: defaultFraction, key)
        self.first = first()
        self.second = second()
    }

    var body: some View {
        GeometryReader { geo in
            let total = axis == .horizontal ? geo.size.width : geo.size.height
            let size = clamp(CGFloat(fraction) * max(0, total - 1), total: total)
            let remaining = max(0, total - size - 1)
            let layout = axis == .horizontal ? AnyLayout(HStackLayout(spacing: 0)) : AnyLayout(VStackLayout(spacing: 0))
            // Resizing changes the layout without replacing the editor or its undo stack.
            layout {
                first.frame(width: axis == .horizontal ? size : nil, height: axis == .vertical ? size : nil)
                handle(total: total, size: size)
                second.frame(width: axis == .horizontal ? remaining : nil, height: axis == .vertical ? remaining : nil)
            }
        }
    }

    private func clamp(_ value: CGFloat, total: CGFloat) -> CGFloat {
        CGFloat(SplitSizing.firstLength(total: Double(total), fraction: total > 1 ? Double(value / (total - 1)) : 0.5,
                                       minFirst: Double(minFirst), minSecond: Double(minSecond)))
    }

    private func handle(total: CGFloat, size: CGFloat) -> some View {
        let available = max(0, total - 1)
        let share = available > 0 ? Double(size / available) : 0
        return Rectangle()
            .fill(Color(nsColor: .separatorColor))
            .frame(width: axis == .horizontal ? 1 : nil, height: axis == .vertical ? 1 : nil)
            .padding(axis == .horizontal ? .horizontal : .vertical, 3)
            .contentShape(Rectangle())
            .padding(axis == .horizontal ? .horizontal : .vertical, -3)
            .pointerStyle(.frameResize(position: axis == .horizontal ? .trailing : .bottom))
            .gesture(DragGesture(minimumDistance: 1, coordinateSpace: .global)
                .onChanged { drag in
                    guard available > 0 else { return }
                    let start = dragStart ?? share
                    if dragStart == nil { dragStart = start }
                    let delta = axis == .horizontal ? drag.translation.width : drag.translation.height
                    fraction = Double(clamp(CGFloat(start) * available + delta, total: total) / available)
                }
                .onEnded { _ in dragStart = nil })
            .accessibilityElement()
            .accessibilityLabel("Divider")
            .accessibilityValue("\(Int(share * 100)) percent")
            .accessibilityAdjustableAction { direction in
                switch direction {
                case .increment: fraction = min(1, share + 0.05)
                case .decrement: fraction = max(0, share - 0.05)
                @unknown default: break
                }
            }
    }
}
