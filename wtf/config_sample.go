package wtf

var sampleConfig = `# Configure WTF application.
app:
  # disable auto-refresh, unless it's explicitly enabled in mods configs
  refreshInterval: -1
  colors:
    background: black
    border:
      focusable: darkslateblue
      focused: orange
      normal: gray
  log:
    level: debug
  grid:
    # configure rows and columns manually by specifying their amount and sizes
    #rows: [13, 10, 4, 13, 13]
    #columns: [35, 20]

    # or specify number of columns, each column's width will be calculated automatically
    numCols: 8
    numRows: 8

# Configure widgets.
# You can configure multiple instances of the same widget type each widget
# should be explicitly enabled, otherwise it's not displayed.
widgets:
  - type: clocks
    enabled: true
    refreshInterval: 15
    sort: "alphabetical"
    position:
      top: 0
      left: 0
      width: 2
      height: 2
    colors:
      rows:
        even: "lightblue"
        odd: "white"
    locations:
      Avignon: "Europe/Paris"
      Barcelona: "Europe/Madrid"
      Dubai: "Asia/Dubai"
      Vancouver: "America/Vancouver"
      Toronto: "America/Toronto"

  - type: system
    enabled: true
    refreshInterval: 3600
    position:
      top: 0
      left: 2
      width: 1
      height: 2

  - type: security
    enabled: true
    refreshInterval: 3600
    position:
      top: 0
      left: 3
      width: 2
      height: 2

  - type: textfile
    enabled: true
    refreshInterval: 1
    filePath: "~/.config/wtf/config.yml"
    format: false
    position:
      top: 0
      left: 5
      width: 3
      height: 7

  - type: logger
    enabled: true
    refreshInterval: 5
    position:
      top: 2
      left: 0
      width: 5
      height: 6
    numLines: 25

  - type: status
    enabled: true
    refreshInterval: 1
    position:
      top: 7
      left: 5
      width: 3
      height: 1
`
