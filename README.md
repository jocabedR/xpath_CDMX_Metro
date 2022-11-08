To run the project the following steps are needed:

    1. Have nodeJS installed.
    2. Locate the project's folder in a terminal and run the following command: 
        * In case you want to test the first phase: 
          go run main.go <input file name> 
          Example: go run main.go Metro_CDMX.kml
        * In case you want to test the second phase: 
          go run main.go <input file name> "<origin>" "<destination>"
          Example: go run main2.go Metro_CDMX.kml "Tacuba" "Tacubaya"