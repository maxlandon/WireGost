package commands

import (
	"github.com/maxlandon/wiregost/client/commands/compiler"
	"github.com/maxlandon/wiregost/client/commands/mainmenu"
	"github.com/maxlandon/wiregost/client/commands/module"
	"github.com/maxlandon/wiregost/client/constants"
	help "github.com/maxlandon/wiregost/client/help/main/core"
)

// BindMain - Binds all commands for the main menu
func BindMain() (err error) {
	// Main help & usage

	// Main unknown handler

	// Register all Commands
	_, err = Main.AddCommand(constants.Compiler, help.CompilerShort, help.CompilerLong, &compiler.Enter{})
	_, err = Main.AddCommand(constants.Exit, "", "", &mainmenu.Exit{})

	_, err = Main.AddCommand(constants.UseModule, "", "", &module.Use{})

	return
}

// BindModule - Binds all commands for the module implant menu
func BindModule() (err error) {
	_, err = Module.AddCommand(constants.Exit, "", "", &mainmenu.Exit{})
	// Main help & usage

	// Main unknown handler

	// Use a further subfunction, that will hide all implant commands
	// that are not compatible with the current implant.
	return
}

// BindGhost - Binds all commands for the ghost implant menu
func BindGhost() {
	// Main help & usage

	// Main unknown handler

	// Use a further subfunction, that will hide all implant commands
	// that are not compatible with the current implant.
}

// BindCompiler - Binds all commands for the compiler menu
func BindCompiler() (err error) {
	// Main help & usage

	// Main unknown handler

	// Register all Commands
	exit, err := Compiler.AddCommand(constants.CompilerExit, help.CompilerExitShort, help.CompilerExitLong, &compiler.Exit{})
	exit.Aliases = []string{constants.CompilerBack, constants.CompilerToMain}

	return
}
