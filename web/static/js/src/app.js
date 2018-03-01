import highcharts from "highcharts"
import $ from "jquery"
import timeago from "timeago"
import moment from "moment"

highcharts.setOptions({
    chart: {
        style: {
            fontFamily:
                '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif',
        },
    },
})

function update() {
    fetch("/api/snapshots")
        .then(function(response) {
            return response.json()
        })
        .then(function(response) {
            if (
                response === null ||
                response.recent === null ||
                response.latest === null
            ) {
                return
            }

            var latest = response.latest
            var recent = response.recent

            if (recent.length === 0) {
                return
            }

            var now = +moment()
            for (let [i, snapshot] of recent.entries()) {
                recent[i].duration = now - +moment(snapshot.timestamp)
            }

            updateModel(latest)

            updateBatteryRemaining(latest)
            updateLoad(latest)
            updateRemainingRuntime(latest)
            updateUtilityVoltage(latest)

            drawBatteryRemainingChart(recent)
            drawLoadChart(recent)
            drawRemainingRuntimeChart(recent)
            drawUtilityVoltageChart(recent)

            $("#last-updated")
                .html("Last updated ")
                .append($("<time />", { datetime: latest.timestamp }))

            $("time").timeago()
        })
}

function humanizeDuration(d) {
    var s = moment.duration(-d).humanize(true)
    s = s[0].toUpperCase() + s.substr(1)
    return s
}

function updateModel(snapshot) {
    $("#model").html("UPS Model: " + snapshot.modelName)
}

function updateBatteryRemaining(snapshot) {
    $("#battery-remaining").html(
        "<h1>" +
            snapshot.batteryRemaining * 100 +
            " %<small>Battery remaining</small></h1>",
    )
}

function updateLoad(snapshot) {
    $("#load").html("<h1>" + snapshot.load + " W<small>Load</small></h1>")
}

function updateRemainingRuntime(snapshot) {
    $("#remaining-runtime").html(
        "<h1>" +
            snapshot.remainingRuntime +
            " min.<small>Remaining runtime</small></h1>",
    )
}

function updateUtilityVoltage(snapshot) {
    $("#utility-voltage").html(
        "<h1>" +
            snapshot.utilityVoltage +
            " V<small>Utility voltage</small></h1>",
    )
}

function getDefaultChartOptions(data) {
    return {
        xAxis: {
            labels: {
                formatter: function() {
                    return humanizeDuration(this.value)
                },
            },
        },

        legend: {
            enabled: false,
        },

        plotOptions: {
            line: {
                marker: {
                    enabled: false,
                },
            },
        },
    }
}

function drawBatteryRemainingChart(snapshots) {
    var data = []
    for (let snapshot of snapshots) {
        data.push([snapshot.duration, snapshot.batteryRemaining])
    }

    highcharts.chart({
        ...getDefaultChartOptions(data),
        chart: {
            renderTo: "battery-remaining-chart",
        },
        title: {
            text: "Battery remaining",
        },

        yAxis: {
            tickInterval: 0.2,
            endOnTick: false,
            labels: {
                formatter: function() {
                    return 100 * this.value + "%"
                },
            },
            min: 0,
            max: 1.1,
        },

        tooltip: {
            formatter: function() {
                return (
                    "<b>" +
                    humanizeDuration(this.x) +
                    ":</b> " +
                    100 * this.y +
                    "%"
                )
            },
        },

        series: [
            {
                name: "Battery remaining",
                data: data,
            },
        ],
    })
}

function drawLoadChart(snapshots) {
    var data = []
    for (let snapshot of snapshots) {
        data.push([snapshot.duration, snapshot.load])
    }

    highcharts.chart({
        ...getDefaultChartOptions(data),
        chart: {
            renderTo: "load-chart",
        },
        title: {
            text: "Load",
        },

        yAxis: {
            endOnTick: true,
            labels: {
                formatter: function() {
                    return this.value + " W"
                },
            },
            min: 0,
            // max: 1 * snapshots[0].batteryCapacity
        },

        tooltip: {
            formatter: function() {
                return (
                    "<b>" + humanizeDuration(this.x) + ":</b> " + this.y + " W"
                )
            },
        },

        series: [
            {
                name: "Load",
                data: data,
            },
        ],
    })
}

function drawRemainingRuntimeChart(snapshots) {
    var data = []
    for (let snapshot of snapshots) {
        data.push([snapshot.duration, snapshot.remainingRuntime])
    }

    highcharts.chart({
        ...getDefaultChartOptions(data),
        chart: {
            renderTo: "remaining-runtime-chart",
        },

        title: {
            text: "Remaining runtime",
        },

        yAxis: {
            endOnTick: true,
            labels: {
                formatter: function() {
                    return this.value + " min"
                },
            },
            min: 0,
        },

        tooltip: {
            formatter: function() {
                return (
                    "<b>" +
                    humanizeDuration(this.x) +
                    ":</b> " +
                    this.y +
                    " min."
                )
            },
        },

        series: [
            {
                name: "Remaining runtime",
                data: data,
            },
        ],
    })
}

function drawUtilityVoltageChart(snapshots) {
    var data = []
    for (let snapshot of snapshots) {
        data.push([snapshot.duration, snapshot.utilityVoltage])
    }

    highcharts.chart({
        ...getDefaultChartOptions(data),
        chart: {
            renderTo: "utility-voltage-chart",
        },

        title: {
            text: "Utility voltage",
        },

        yAxis: {
            endOnTick: true,
            labels: {
                formatter: function() {
                    return this.value + " V"
                },
            },
            min: 0,
        },

        tooltip: {
            formatter: function() {
                return (
                    "<b>" + humanizeDuration(this.x) + ":</b> " + this.y + " V"
                )
            },
        },

        series: [
            {
                name: "Utility voltage",
                data: data,
            },
        ],
    })
}

function setMonitoringStartTime() {
    $("#started")
        .append("Monitoring started ")
        .append($("<time />", { datetime: StartTime }))

    $("time").timeago()
}

function setRevision() {
    var revision = window.REVISION

    $("#revision")
        .append(
            $("<a />", {
                href:
                    "https://github.com/jimmysawczuk/power-monitor/commit/" +
                    revision.hex.full,
            }).html("rev. " + revision.hex.short),
        )
        .append(" &middot; ")
        .append($("<time />", { datetime: revision.commit_date.iso8601 }))

    $("time").timeago()
}

document.addEventListener("DOMContentLoaded", function() {
    window.setInterval(update, 30 * 1000)
    update()

    setMonitoringStartTime()
    setRevision()
})
