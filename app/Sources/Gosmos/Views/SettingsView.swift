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
        .frame(width: 560)
        .scenePadding()
    }
}

struct GeneralSettings: View {
    @Environment(AppModel.self) private var app

    var body: some View {
        Form {
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
                            Text(library.error?.title ?? "\(library.catalog.lessonCount) lessons, \(app.done[library.id]?.count ?? 0) done")
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
    }
}

struct EditorSettings: View {
    @Environment(AppModel.self) private var app
    private static let sample = "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, gopher\")\n}\n"

    var body: some View {
        Form {
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
    }
}

struct GoSettings: View {
    @Environment(AppModel.self) private var app
    @ViewState private var path = ""

    var body: some View {
        Form {
            Section {
                TextField("Go command", text: $path, prompt: Text("Found automatically"))
                    .font(.system(.body, design: .monospaced))
                    .onSubmit { app.rediscoverToolchain(goPath: path) }
                HStack {
                    Button("Choose…") {
                        let panel = NSOpenPanel()
                        panel.canChooseFiles = true
                        panel.canChooseDirectories = false
                        panel.message = "Choose the go command, for example /usr/local/go/bin/go"
                        if panel.runModal() == .OK, let url = panel.url {
                            path = url.path
                            app.rediscoverToolchain(goPath: path)
                        }
                    }
                    Button("Find Automatically") {
                        path = ""
                        app.rediscoverToolchain(goPath: nil)
                    }
                }
            } header: {
                Text("Toolchain")
            } footer: {
                Text("Gosmos runs your code with this Go. Left empty, it uses the go your login shell finds.")
            }
            Section("Detected") {
                if let t = app.toolchain {
                    if let go = t.go {
                        LabeledContent("Version", value: t.goVersion ?? "unknown")
                        LabeledContent("Path") { Text(go).font(.system(.body, design: .monospaced)).textSelection(.enabled) }
                    } else {
                        GoMissingView()
                    }
                } else {
                    HStack { ProgressView().controlSize(.small); Text("Looking for Go…") }
                }
            }
        }
        .formStyle(.grouped)
        .onAppear { path = app.preferences.goPath ?? "" }
    }
}

/// How to install Go when it is missing.
struct GoMissingView: View {
    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            Label("Go is not installed", systemImage: "exclamationmark.triangle.fill")
                .symbolRenderingMode(.multicolor)
                .fontWeight(.medium)
            Text("Install it from go.dev/dl, or with Homebrew: brew install go. Then choose Find Automatically.")
                .foregroundStyle(.secondary)
            if let download = URL(string: "https://go.dev/dl/") {
                Link("Download Go", destination: download)
            }
        }
    }
}

