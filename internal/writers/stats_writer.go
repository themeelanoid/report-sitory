package writers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/themeelanoid/report-sitory/internal/stats"
)

type StatsWriter interface {
	WriteStats(statistics []stats.StatLine) error
}

func WriteStatistics(statistics []stats.StatLine, format string) error {
	var writer StatsWriter
	switch format {
	case "tabular":
		writer = NewTabularWriter(os.Stdout)
	case "csv":
		writer = NewCSVWriter(os.Stdout)
	case "json":
		writer = NewJSONWriter(os.Stdout)
	case "json-lines":
		writer = NewJSONLinesWriter(os.Stdout)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
	err := writer.WriteStats(statistics)
	return err
}

type CSVWriter struct {
	writer *csv.Writer
}

type JSONLinesWriter struct {
	writer io.Writer
}

type JSONWriter struct {
	writer io.Writer
}

type TabularWriter struct {
	writer *tabwriter.Writer
}

func NewCSVWriter(w io.Writer) *CSVWriter {
	_, _ = w.Write([]byte("Name,Lines,Commits,Files\n"))
	return &CSVWriter{writer: csv.NewWriter(w)}
}

func NewJSONLinesWriter(w io.Writer) *JSONLinesWriter {
	return &JSONLinesWriter{writer: w}
}

func NewJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{writer: w}
}

func NewTabularWriter(w io.Writer) *TabularWriter {
	res := &TabularWriter{writer: tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)}
	_, _ = fmt.Fprintln(res.writer, "Name\tLines\tCommits\tFiles")
	return res
}

func (w *CSVWriter) WriteStats(stats []stats.StatLine) error {
	for _, line := range stats {
		lines := strconv.Itoa(line.Lines)
		commits := strconv.Itoa(line.Commits)
		files := strconv.Itoa(line.Files)
		_ = w.writer.Write([]string{line.Author, lines, commits, files})
	}
	w.writer.Flush()
	return nil
}

func (w *JSONLinesWriter) WriteStats(stats []stats.StatLine) error {
	for _, line := range stats {
		bytes, err := json.Marshal(line)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w.writer, "%s\n", string(bytes))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *JSONWriter) WriteStats(frames []stats.StatLine) error {
	bytes, err := json.Marshal(frames)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(bytes)
	return err
}

func (w *TabularWriter) WriteStats(stats []stats.StatLine) error {
	for _, line := range stats {
		_, err := fmt.Fprintf(w.writer, "%s\t%d\t%d\t%d\n", line.Author, line.Lines, line.Commits, line.Files)
		if err != nil {
			return err
		}
	}
	return w.writer.Flush()
}
