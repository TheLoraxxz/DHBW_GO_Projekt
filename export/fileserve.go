package export

/*
Um später einzubinden:
http.HandleFunc("/", export.WrapperAuth(export.AuthenticatorFunc(export.CheckUserValid), export.Handler))
log.Fatalln(http.ListenAndServe(":8080", nil))
*/

import (
	"DHBW_GO_Projekt/authentifizierung"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Handler ToDo relativePfade + Parametrisierung vom Handler
func Handler(w http.ResponseWriter, r *http.Request) { //für Download nötiger Handle

	cookie, err := r.Cookie("Download-Kalender")
	if err != nil {
		return
	}
	_, user := authentifizierung.CheckCookie(&cookie.Value)

	err = serveFile(w, r, "C:\\Users\\chris\\Documents\\GitHub\\DHBW_GO_Projekt\\export\\export\\"+"mik"+".ical")
	if err != nil {
		return
	}
}

func serveFile(writer http.ResponseWriter, request *http.Request, filePath string) (err error) { //stellt sicher, dass der Download pausiert werden kann
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s""`, fileInfo.Name()))
	writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	writer.Header().Set("Accept-Ranges", "bytes")
	requestRange := request.Header.Get("range")
	if requestRange == "" {
		writer.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))
		_, err := file.Seek(0, 0)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
		return nil
	}
	requestRange = requestRange[6:]
	splitRange := strings.Split(requestRange, "-")
	if len(splitRange) != 2 {
		return fmt.Errorf("invalid values for header 'Range'")
	}
	begin, err := strconv.ParseInt(splitRange[0], 10, 64)
	if err != nil {
		return err
	}
	end, err := strconv.ParseInt(splitRange[1], 10, 64)
	if err != nil {
		return err
	}
	if begin > fileInfo.Size() || end > fileInfo.Size() {
		return fmt.Errorf("range out of bounds for file")
	}
	if begin >= end {
		return fmt.Errorf("range begin cannot be bigger than range end")
	}
	writer.Header().Set("Content-Length", strconv.FormatInt(end-begin+1, 10))
	writer.Header().Set("Content-Range",
		fmt.Sprintf("bytes %d-%d/%d", begin, end, fileInfo.Size()))
	writer.WriteHeader(http.StatusPartialContent)
	_, err = file.Seek(begin, 0)
	if err != nil {
		return err
	}
	_, err = io.CopyN(writer, file, end-begin)
	if err != nil {
		return err
	}
	return nil
}
