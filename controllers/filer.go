package controllers

import (
    "fmt"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "encoding/json"
    "log"
    "github.com/cschar/pp-gocrawler/imageprocessing"
)


// UploadFile uploads a file to the server
func UploadFile(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    file, handle, err := r.FormFile("file")
    if err != nil {
        fmt.Fprintf(w, "%v", err)
        return
    }
    defer file.Close()

    mimeType := handle.Header.Get("Content-Type")
    switch mimeType {
    case "image/jpeg":
        saveFile(w, file, handle)
    //case "image/png":
    //    saveFile(w, file, handle)
    default:
        jsonResponse(w, http.StatusBadRequest, "The format file is not valid.")
    }
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
    data, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Fprintf(w, "%v", err)
        return
    }

    filepath := "./uploadfiles/"+handle.Filename
    err = ioutil.WriteFile(filepath, data, 0666)
    if err != nil {
        fmt.Fprintf(w, "%v", err)
        return
    }

    //mix image
    imageprocessing.MakeImageFromSlices(filepath)



    jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    fmt.Fprint(w, message)
}


type Profile struct {
  Name    string
  Files []string
}

func MixedImages(w http.ResponseWriter, r *http.Request) {
  //profile := Profile{"ImageMixes", []string{"output/eyemazestyle.png", "output/snowymandala.png"}}

    files, err := ioutil.ReadDir("./public/output")  // relative to main.go
    if err != nil {
        log.Fatal(err)
    }

    var s []string
    for _, f := range files {
        s = append(s, f.Name())
    }
    profile2 := Profile{"ImageMixes", s}

  js, err := json.Marshal(profile2)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}


func InputImages(w http.ResponseWriter, r *http.Request) {
  //profile := Profile{"ImageMixes", []string{"output/eyemazestyle.png", "output/snowymandala.png"}}

    files, err := ioutil.ReadDir("./public/input")  // relative to main.go
    if err != nil {
        log.Fatal(err)
    }

    var s []string
    for _, f := range files {
        s = append(s, f.Name())
    }
    profile2 := Profile{"ImageInputs", s}

  js, err := json.Marshal(profile2)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}