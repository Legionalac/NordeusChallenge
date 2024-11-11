package main

import (
	"fmt"
	"net/http"
	"nordeuschallenge/components"
	"nordeuschallenge/libs"
	"nordeuschallenge/middlewares"
	"slices"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Position struct {
    x int
    y int
}

var logger = libs.GetLogger();
var islandMap [][]int;
var islandMax = 0;

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx.Response().Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
	ctx.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "true")
	ctx.Response().Header().Set(echo.HeaderAccessControlAllowHeaders , "true")
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func exploreIsland(islandMap [][]int,visited [][]bool, x int, y int) int{
    stack := []Position{
        {x,y},
    }

    sum:= 0
    count:= 0
    for len(stack) > 0{
        currentPosition := stack[len(stack) - 1]
        stack = stack[:len(stack) - 1]
        count++

        sum += islandMap[currentPosition.x][currentPosition.y]
        visited[currentPosition.x][currentPosition.y] = true
        //RIGHT
        if currentPosition.x + 1 < 30 &&
        islandMap[currentPosition.x+1][currentPosition.y] > 0 &&
        !visited[currentPosition.x+1][currentPosition.y]{
            visited[currentPosition.x + 1][currentPosition.y] = true
            stack = append(stack, Position{currentPosition.x + 1, currentPosition.y})
        }
        //LEFT
        if currentPosition.x - 1 >= 0 &&
        islandMap[currentPosition.x - 1][currentPosition.y] > 0 &&
        !visited[currentPosition.x-1][currentPosition.y]{
            visited[currentPosition.x - 1][currentPosition.y] = true
            stack = append(stack, Position{currentPosition.x - 1, currentPosition.y})
        }
        //UP
        if currentPosition.y + 1 < 30 &&
        islandMap[currentPosition.x][currentPosition.y + 1] > 0 &&
        !visited[currentPosition.x][currentPosition.y + 1]{
            visited[currentPosition.x][currentPosition.y + 1] = true
            stack = append(stack, Position{currentPosition.x, currentPosition.y + 1})
        }
        //DOWN
        if currentPosition.y - 1 >= 0 &&
        islandMap[currentPosition.x][currentPosition.y - 1] > 0 &&
        !visited[currentPosition.x][currentPosition.y - 1]{
            visited[currentPosition.x][currentPosition.y - 1] = true
            stack = append(stack, Position{currentPosition.x, currentPosition.y - 1})
        }
    }

    return sum/count


}


func findHighestIsland(islandMap [][]int) int{

    visited := libs.InitVisitedMatrix()

	var islandAverages []int

    for i := range visited{
        for j := range visited[i]{
			if islandMap[i][j] > 0 && !visited[i][j] {
                average := exploreIsland(islandMap, visited, i, j)
                islandAverages = append(islandAverages, average)
			}
        }
    }
    logger.Info().Msg("Island averages are " + fmt.Sprint(islandAverages))
    return slices.Max(islandAverages)
}

func main(){
    // Echo instance
    e := echo.New()

    // Middleware
    middlewares.LoggerMiddleware(e, logger)
    e.Use(middleware.Recover())

    // Routes
    e.GET("/", func(c echo.Context) error {
        islandMap = libs.GetMatrix()
        islandMax = findHighestIsland(islandMap)
        return Render(c, http.StatusOK, components.Index(islandMap))
    })

    e.POST("/attempt", func(c echo.Context) error {
        visited := libs.InitVisitedMatrix()

        x, err := strconv.Atoi(c.FormValue("rowId"))

        if err!=nil{
            logger.Err(err).Msg("Error while parsing string " + err.Error())
        }

        y, err := strconv.Atoi(c.FormValue("colId"))

        if err!=nil{
            logger.Err(err).Msg("Error while parsing string " + err.Error())
        }
        selectedMax := exploreIsland(islandMap, visited, x,y)
        logger.Info().Msg("ISLAND MAX" + strconv.Itoa(islandMax) + " SELECTED MAX " + strconv.Itoa(selectedMax) )
        if(selectedMax == islandMax){
            return Render(c,http.StatusOK, components.Victory())
        } else{
            return Render(c,http.StatusOK, components.Defeat())
        }
    })

    e.Static("/static/*", "public/static")
    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}

