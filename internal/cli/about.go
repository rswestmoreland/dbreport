package cli

import (
	"fmt"
	"io"

	"github.com/rswestmoreland/dbreport/internal/version"
)

func PrintVersion(w io.Writer) {
	fmt.Fprintf(w, "%s\n", version.Short())
	fmt.Fprintf(w, "Build date: %s\n", version.Date)
	fmt.Fprintf(w, "Commit: %s\n", version.Commit)
}

func PrintAbout(w io.Writer) {
	fmt.Fprintf(w, "%s %s\n\n", version.AppName, version.Version)
	fmt.Fprintf(w, "%s.\n\n", version.Description)

	fmt.Fprintln(w, "Author:")
	fmt.Fprintf(w, "  %s\n", version.Author)
	fmt.Fprintf(w, "  %s\n\n", version.Email)

	fmt.Fprintln(w, "License:")
	fmt.Fprintf(w, "  %s License\n\n", version.License)

	fmt.Fprintln(w, "Copyright:")
	fmt.Fprintf(w, "  %s\n\n", version.Copyright)

	fmt.Fprintln(w, "Build:")
	fmt.Fprintf(w, "  Version: %s\n", version.Version)
	fmt.Fprintf(w, "  Date: %s\n", version.Date)
	fmt.Fprintf(w, "  Commit: %s\n", version.Commit)
}
