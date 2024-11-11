package libs

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var logger = GetLogger();

func GetMatrix() [][] int{
    start := time.Now()
    resp, err := http.Get("https://jobfair.nordeus.com/jf24-fullstack-challenge/test")
    if err != nil {
        logger.Error().Msg(err.Error());
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        logger.Error().Msg(err.Error())
    }

    stringResponse := string(body)
    elapsed := time.Since(start)

    logger.Info().Msg("Api call took " + elapsed.String())
    matrix, err := convertStringToMatrix(stringResponse)

    if err!=nil {
        logger.Error().Msg(err.Error())
    }

    return matrix
}

func convertStringToMatrix(stringValue string) ([][] int, error){
    lines := strings.Split(stringValue, "\n")

    var matrix [][]int

    for _, line := range lines {
        strNumbers := strings.Fields(line)

        var row []int

        for _, strNum := range strNumbers {
            num, err := strconv.Atoi(strNum)
            if err != nil {
                return nil, errors.New("Error converting string to integer")
            }
            row = append(row, num)
        }

        // Append the row to the matrix
        matrix = append(matrix, row)
    }

    return matrix, nil
}

func InitVisitedMatrix() [][]bool {

    visited := make([][]bool, 30)

    for i:= range visited{
        visited[i] = make([]bool, 30)
    }

    for i := range visited{
        for j := range visited[i]{
            visited[i][j] = false
        }
    }

    return visited
}
