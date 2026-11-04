import Testing
@testable import GosmosKit

@Suite struct SplitSizingTests {
    @Test func panesStayInsideANarrowWindow() {
        for width in [0.0, 1, 200, 440, 579, 580, 600, 1000] {
            for fraction in [-1.0, 0, 0.3, 0.9, 2, .nan, .infinity] {
                let first = SplitSizing.firstLength(total: width, fraction: fraction, minFirst: 260, minSecond: 320)
                let available = max(0, width - 1)
                #expect(first.isFinite)
                #expect(first >= 0 && first <= available)
                #expect(available - first >= 0)
                if available >= 580 {
                    #expect(first >= 260)
                    #expect(available - first >= 320)
                }
            }
        }
    }

    @Test func narrowPanesShareTheAvailableSpace() {
        let first = SplitSizing.firstLength(total: 291, fraction: 0.9, minFirst: 260, minSecond: 320)
        #expect(first == 130)
        #expect(290 - first == 160)
    }

    @Test func aCorruptSavedFractionStillProducesAUsableSplit() {
        #expect(SplitSizing.firstLength(total: 1001, fraction: .nan, minFirst: 200, minSecond: 380) == 500)
        #expect(SplitSizing.firstLength(total: 1001, fraction: -2, minFirst: 200, minSecond: 380) == 200)
        #expect(SplitSizing.firstLength(total: 1001, fraction: 2, minFirst: 200, minSecond: 380) == 620)
    }
}
