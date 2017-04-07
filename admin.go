package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func getRootHandler(state *string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templateDef := `
<html>
    <head>
        <style>
            body {
                text-align: center;
                font-family: "Verdana";
            }

            .main {
                width: 90%;
                margin: auto;
                margin-top: 20px;
            }

            .help {
                text-align: left;
                padding: 10px;
            }

            .state {
                padding: 10px;
                color: white;
            }

            .state-down {
                background:red;
            }
            
            .border-state-down {
                border: 1px solid red;
            }

            .state-up, .state-ready {
                background: green;
            }

            .border-state-up, .border-state-ready {
                border: 1px solid green;
            }

            .state-maint {
                background: darkorange;
            }

            .border-state-maint {
                border: 1px solid darkorange;
            }

            .state-drain {
                background: darkblue;
            }

            .border-state-drain {
                border: 1px solid darkblue;
            }
        </style>
    </head>
    <body>
<svg width="100px" height="100px" viewBox="0 0 209 229" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <!-- Generator: Sketch 42 (36781) - http://www.bohemiancoding.com/sketch -->
    <title>Logo_Vertical</title>
    <desc>Created with Sketch.</desc>
    <defs></defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="iPad-Pro-Portrait" transform="translate(-827.000000, -998.000000)">
            <g id="Logo_Vertical" transform="translate(749.000000, 920.000000)">
                <polygon id="Fill-2" fill="#FF585D" points="262.2728 237.992 262.2728 158.228 182.5088 158.228 182.5088 158.228 182.5088 158.228 182.5088 78.47 102.7508 78.47 102.7508 237.992"></polygon>
                <polygon id="Fill-3" fill="#2CD5C4" points="262.2728 78.4676 206.4368 78.4676 206.4368 134.3036 262.2728 134.3036"></polygon>
                <path d="M198.1262,299.1758 C197.1542,299.4278 195.9122,299.5838 194.8202,299.5838 C189.7022,299.5838 186.1322,295.9778 186.1322,290.8118 C186.1322,285.6758 189.8222,282.0878 195.1022,282.0878 C198.9062,282.0878 201.7202,283.1978 203.5082,284.0318 L204.4082,284.4518 L204.4082,277.3598 L204.0422,277.1918 C201.1382,275.8478 197.9942,275.1638 194.7002,275.1638 C185.5502,275.1638 178.6502,281.8898 178.6502,290.8118 C178.6502,299.8838 185.4482,306.4658 194.8202,306.4658 C198.3422,306.4658 202.2242,305.5898 204.9482,304.1858 L205.2902,304.0118 L205.2902,289.2698 L198.1262,289.2698 L198.1262,299.1758 Z M233.4782,305.9018 L240.9182,305.9018 L240.9182,294.1718 L250.5602,294.1718 L250.5602,287.8118 L240.9182,287.8118 L240.9182,282.1658 L252.3662,282.1658 L252.3662,275.7278 L233.4782,275.7278 L233.4782,305.9018 Z M277.7402,275.7278 L271.4822,286.9118 L265.2242,275.7278 L256.8002,275.7278 L267.7562,294.2378 L267.7562,305.9018 L275.2022,305.9018 L275.2022,294.2378 L286.2002,275.7278 L277.7402,275.7278 Z M215.3642,305.9018 L222.8102,305.9018 L222.8102,275.7278 L215.3642,275.7278 L215.3642,305.9018 Z M114.0662,290.8118 C114.0662,285.8378 117.6302,282.0878 122.3522,282.0878 C127.0802,282.0878 130.6382,285.8378 130.6382,290.8118 C130.6382,295.7918 127.0802,299.5418 122.3522,299.5418 C117.6302,299.5418 114.0662,295.7918 114.0662,290.8118 Z M106.5842,290.8118 C106.5842,299.7338 113.3642,306.4658 122.3522,306.4658 C131.3462,306.4658 138.1262,299.7338 138.1262,290.8118 C138.1262,282.0398 131.1962,275.1638 122.3522,275.1638 C113.5082,275.1638 106.5842,282.0398 106.5842,290.8118 Z M152.5142,299.4638 L152.5142,282.1658 L154.8602,282.1658 C160.4162,282.1658 164.1482,285.6398 164.1482,290.8118 C164.1482,295.9838 160.4462,299.4638 154.9382,299.4638 L152.5142,299.4638 Z M145.0742,275.7278 L145.0742,305.9018 L155.0222,305.9018 C164.9582,305.9018 171.6362,299.8358 171.6362,290.8118 C171.6362,281.7878 164.9222,275.7278 154.9382,275.7278 L145.0742,275.7278 Z M86.2622,275.7278 L78.8222,275.7278 L78.8222,305.9018 L101.4602,305.9018 L101.4602,299.4638 L86.2622,299.4638 L86.2622,275.7278 Z" id="Fill-4" fill="#333333"></path>
            </g>
        </g>
    </g>
</svg>
        <div class="main border-state-{{.State}}">
            <div class="state state-{{.State}}">
                <b>
                    Current Agent State: {{.State}}</span>
                </b>
            </div>
            <div>
                <form action="/" method="post">
                    <b>
                        New state :
                    </b>
                    <select name="state">
                        <option value="up">Up</option>
                        <option value="ready">Ready</option>
                        <option value="maint">Maintenance</option>
                        <option value="drain">Draining</option>
                        <option value="down">Down</option>
                    </select>
                    <input type="submit" value="Commit" />
                </form>
            </div>
            <div>
                <p class="help">
                    <b>UP : </b> sets back the server's operating state as UP if health checks
    also report that the service is accessible.
                </p>
                <p class="help">
                    <b>READY : </b> This will turn the server's administrative state to the
        READY mode, thus cancelling any DRAIN or MAINT state
                </p>
                <p class="help">
                    <b>DRAIN: </b> This will turn the server's administrative state to the
        DRAIN mode, thus it will not accept any new connections other than those
        that are accepted via persistence.
                </p>
                <p class="help">
                    <b>MAINT: </b> This will turn the server's administrative state to the
        MAINT mode, thus it will not accept any new connections at all, and health
        checks will be stopped.
                </p>
                <p class="help">
                    <b>DOWN: </b> This will mark the server's
        operating state as DOWN
                </p>
            </div>
        </div>
    </body>
</html>
`
		if r.Method == "POST" {
			newState := r.PostFormValue("state")
			if isValidState(newState) {
				log.Printf("New state set to %s by %s\n", newState, r.RemoteAddr)
				*state = newState
			}
		}

		data := struct {
			State *string
		}{State: state}
		tpl, _ := template.New("root").Parse(templateDef)
		tpl.Execute(w, data)
	}
}

func startAdministrativeInterface(port int, state *string) {
	http.HandleFunc("/", getRootHandler(state))
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), nil)
}
