import AppKit
import CodeEditor
import GosmosKit
import SwiftUI

struct SettingsView: View {
    /// The pane last shown, so Settings opens where you left it.
    @AppStorage("SettingsTab") private var tab = "general"

    var body: some View {
        TabView(selection: $tab) {
            Tab("General", systemImage: "gearshape", value: "general") { GeneralSettings() }
            Tab("Editor", systemImage: "character.cursor.ibeam", value: "editor") { EditorSettings() }
            Tab("Go", systemImage: "chevron.left.forwardslash.chevron.right", value: "go") { GoSettings() }
        }
        .frame(width: 610, height: 640)
        .scenePadding()
        .background(Theme.canvas)
    }
}

struct GeneralSettings: View {
    @Environment(AppModel.self) private var app
    @ViewState private var showsAcknowledgements = false

    var body: some View {
        Form {
            Section {
                MissionHeading(eyebrow: "MAKE YOURSELF AT HOME", title: "Your learning space.", subtitle: "Your courses, your progress, your pace.", symbol: "slider.horizontal.3")
                    .padding(.vertical, 8)
            }
            Section("Your journey") {
                LabeledContent("Explorer rank", value: "Level \(app.learningProgress.level) · \(app.learningProgress.rank)")
                LabeledContent("Experience", value: "\(app.learningProgress.xp.formatted()) XP")
                Text("Every completed lesson contributes \(LearningProgress.lessonXP) XP. Your saved progress powers your rank and badges, including lessons you already finished.")
                    .font(.caption).foregroundStyle(Theme.secondaryText)
            }
            Section {
                LabeledContent("The language and the standard library") {
                    Button("Acknowledgements…") { showsAcknowledgements = true }
                }
            } header: {
                Text("Go Documentation")
            } footer: {
                Text("The specification, the guides and the standard library's documentation are built into Gosmos. The acknowledgements credit their authors and give their licenses.")
            }
            Section {
                ForEach(app.shelf) { library in
                    LabeledContent {
                        if library.isBuiltIn {
                            Text("Built in").foregroundStyle(Theme.secondaryText)
                        } else {
                            Button("Show in Finder") { NSWorkspace.shared.activateFileViewerSelecting([library.url]) }
                                .help(library.url.path)
                        }
                    } label: {
                        VStack(alignment: .leading, spacing: 2) {
                            Text(library.title)
                            // Done counts only lessons the library still has, as XP above does.
                            Text(library.error?.title ?? "\(library.catalog.lessonCount) lessons, \((app.done[library.id] ?? []).filter { library.catalog.summary($0) != nil }.count) done")
                                .font(.caption)
                                .foregroundStyle(Theme.secondaryText)
                        }
                    }
                }
                Button("Open Library…") { app.chooseLibrary() }
                    .disabled(app.learner == nil || app.isAddingLibrary)
            } header: {
                Text("Libraries")
            } footer: {
                Text("The courses built into Gosmos, then the library files you opened. Gosmos only reads them; your progress stays on this Mac.")
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
        .sheet(isPresented: $showsAcknowledgements) { AcknowledgementsSheet(reference: app.docs?.reference) }
    }
}

struct EditorSettings: View {
    @Environment(AppModel.self) private var app
    private static let sample = "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, gopher\")\n}\n"

    var body: some View {
        Form {
            Section {
                MissionHeading(eyebrow: "BUILT FOR YOUR FLOW", title: "Room to create.", subtitle: "Make the editor feel like yours.", symbol: "character.cursor.ibeam", color: Theme.violet)
                    .padding(.vertical, 8)
            }
            Section("Editor") {
                LabeledContent("Font size") {
                    Stepper(value: Binding(get: { app.editorFontSize }, set: { app.setEditorFontSize($0) }), in: 9...28, step: 1) {
                        Text("\(Int(app.editorFontSize)) pt").monospacedDigit()
                    }
                }
                Toggle("Line numbers", isOn: Binding(get: { app.showLineNumbers }, set: { app.setShowLineNumbers($0) }))
                CodeEditorView(
                    text: .constant(Self.sample),
                    highlighter: GoLexer(),
                    fontSize: app.editorFontSize,
                    isEditable: false,
                    showsLineNumbers: app.showLineNumbers,
                    accessibilityLabel: "Preview"
                )
                .frame(height: 150)
                .clipShape(RoundedRectangle(cornerRadius: 8, style: .continuous))
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
    }
}

struct GoSettings: View {
    @Environment(AppModel.self) private var app

    var body: some View {
        Form {
            Section {
                MissionHeading(eyebrow: "READY FOR LIFTOFF", title: "Your engine for Go.", subtitle: "Everything you need to build, run, and explore.", symbol: "terminal", color: Theme.accent)
                    .padding(.vertical, 8)
            }
            Section {
                if let t = app.toolchain {
                    if let name = t.name {
                        LabeledContent("Runs your code", value: name)
                        if !t.isBuiltIn, let go = t.go {
                            LabeledContent("Path") { Text(go).font(.system(.body, design: .monospaced)).textSelection(.enabled) }
                        }
                    } else {
                        Label {
                            Text(t.problem ?? "There is no Go to run your code with.")
                        } icon: {
                            Image(systemName: "exclamationmark.triangle.fill").symbolRenderingMode(.multicolor)
                        }
                    }
                } else {
                    HStack { ProgressView().controlSize(.small); Text("Checking Go…") }
                }
            } header: {
                Text("Go")
            } footer: {
                Text("Gosmos builds, runs, tests and formats your code with the Go built into it, offline. Nothing needs installing, and nothing is downloaded.")
            }
            Section {
                HStack {
                    Button("Choose…") {
                        let panel = NSOpenPanel()
                        panel.canChooseFiles = true
                        panel.canChooseDirectories = false
                        panel.message = "Choose the go command of another Go, for example /usr/local/go/bin/go"
                        if panel.runModal() == .OK, let url = panel.url {
                            app.rediscoverToolchain(goPath: url.path)
                        }
                    }
                    Button("Use Built-in Go") { app.rediscoverToolchain(goPath: nil) }
                        .disabled(app.toolchain?.isBuiltIn != false)
                }
            } header: {
                Text("Another Go")
            } footer: {
                Text("For advanced use: run your code with another Go on this Mac, such as a newer release. It runs offline too.")
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
    }
}
