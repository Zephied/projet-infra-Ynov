# STEPS FOR RUNNING THE APPLICATION

* **Step1:** make sure you have Golang installed and working well on your code editor, if not download and install Go for your OS. Once this is done
 - install golang via this link.  https://go.dev/dl/

make sure to run this code on the terminal to allow golang to recognise the application as a module, first navigate to the project directory and run e.g `gui-chatbot`

```bash
go mod init gui-chatbot
```

* **Step 2:** Create a python virtual environment and make sure to install streamlit if not yet installed
  
* **Step 3:** Once step 1 and 2 are done! Run the `server.go` file on the terminal with the substeps;

  * build the script

    ```bash
        go build
    ```

  * GOTO the `Makefile` file, copy and paste the code to the terminal to **create the environment variable for the api key**

  * run the server

  ```bash
    go run server.go
  ```

* **step 5:** Once you have the go server running on a terminal, create a seperate terminal to run the streamlit application for the `app.py`

  ```bash
       sudo apt install python3-venv
    ```

      ```bash
       python3 -m venv path/to/venv
    ```

      ```bash
      source path/to/venv/bin/activate
    ```

     ```bash
      pip install streamlit
    ```

    ```bash
        streamlit run app.py --server.headless true
    ```

* **step 6:** Also, create another seperate terminal to run the streamlit application for the `spotify-app.py`

    ```bash
        streamlit run spotify-app.py --server.headless true
    ```

* **step 6:** To view the Homepage of the API.

    ```bash
        navigate to `http://127.0.0.1:5500/`

        or open the `index.html` in a browser.
    ```
