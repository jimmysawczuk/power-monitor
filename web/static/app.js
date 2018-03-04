import highcharts from 'highcharts'
import $ from 'jquery-slim'
import timeago from 'timeago.js'
import moment from 'moment'
import 'whatwg-fetch'
import fontawesome from '@fortawesome/fontawesome'

import './style.less'

fontawesome.library.add({
	prefix:
		'fa',
	iconName:
		'battery-bolt',
	icon: [
		640,
		512,
		[],
		'f376',
		'M64 352h178.778l-14.173 64H48c-26.51 0-48-21.49-48-48V144c0-26.51 21.49-48 48-48h115.944l-7.663 64H64v192zm364.778-160h-92.321l27.694-133.589C367.4 45.087 358.205 32 345.6 32H236.8c-9.623 0-17.76 7.792-19.031 18.225L192.171 264c-1.535 12.59 7.432 23.775 19.031 23.775h94.961l-36.847 166.382C266.44 467.443 275.728 480 287.993 480c6.68 0 13.101-3.827 16.623-10.481l140.778-245.997C452.79 209.55 443.564 192 428.778 192zM616 160h-8v-16c0-26.51-21.49-48-48-48H405.38l-9.951 48h33.349c16.112 0 31.233 5.762 43.115 16H544v64h32v64h-32v64H427.174l-36.626 64H560c26.51 0 48-21.49 48-48v-16h8c13.255 0 24-10.745 24-24V184c0-13.255-10.745-24-24-24z',
	],
})

highcharts.setOptions({
	chart: {
		style: {
			fontFamily:
				'-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif',
		},
	},
})

function humanizeDuration(d) {
	return `${moment
		.duration(
			d,
			'milliseconds',
		)
		.humanize()} ago`
}

function updateModel(snapshot) {
	$('#model').html(`UPS Model: ${
		snapshot.modelName
	}`)
}

function updateBatteryRemaining(snapshot) {
	$('#battery-remaining').html(`<h1>${snapshot.batteryRemaining *
			100} %<small>Battery remaining</small></h1>`)
}

function updateLoad(snapshot) {
	$('#load').html(`<h1>${
		snapshot.load
	} W<small>Load</small></h1>`)
}

function updateRemainingRuntime(snapshot) {
	$('#remaining-runtime').html(`<h1>${
		snapshot.remainingRuntime
	} min.<small>Remaining runtime</small></h1>`)
}

function updateUtilityVoltage(snapshot) {
	$('#utility-voltage').html(`<h1>${
		snapshot.utilityVoltage
	} V<small>Utility voltage</small></h1>`)
}

function getDefaultChartOptions() {
	return {
		xAxis: {
			labels: {
				formatter() {
					return humanizeDuration(this
						.value)
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
	const data = []
	snapshots.forEach((snapshot) => {
		data.push([
			snapshot.duration,
			snapshot.batteryRemaining,
		])
	})

	highcharts.chart({
		...getDefaultChartOptions(),
		chart: {
			renderTo:
					'battery-remaining-chart',
		},
		title: {
			text:
					'Battery remaining',
		},

		yAxis: {
			tickInterval: 0.2,
			endOnTick: false,
			labels: {
				formatter() {
					return `${100 *
							this
								.value}%`
				},
			},
			min: 0,
			max: 1.1,
		},

		tooltip: {
			formatter() {
				return `<b>${humanizeDuration(this
					.x)}:</b> ${100 *
						this
							.y}%`
			},
		},

		series: [
			{
				name:
						'Battery remaining',
				data,
			},
		],
	})
}

function drawLoadChart(snapshots) {
	const data = []
	snapshots.forEach((snapshot) => {
		data.push([
			snapshot.duration,
			snapshot.load,
		])
	})

	highcharts.chart({
		...getDefaultChartOptions(),
		chart: {
			renderTo:
					'load-chart',
		},
		title: {
			text:
					'Load',
		},
		yAxis: {
			endOnTick: true,
			labels: {
				formatter() {
					return `${
						this
							.value
					} W`
				},
			},
			min: 0,
			// max: 1 * snapshots[0].batteryCapacity
		},

		tooltip: {
			formatter() {
				return `<b>${humanizeDuration(this
					.x)}:</b> ${
					this
						.y
				} W`
			},
		},

		series: [
			{
				name:
						'Load',
				data,
			},
		],
	})
}

function drawRemainingRuntimeChart(snapshots) {
	const data = []
	snapshots.forEach((snapshot) => {
		data.push([
			snapshot.duration,
			snapshot.remainingRuntime,
		])
	})

	highcharts.chart({
		...getDefaultChartOptions(),
		chart: {
			renderTo:
					'remaining-runtime-chart',
		},

		title: {
			text:
					'Remaining runtime',
		},

		yAxis: {
			endOnTick: true,
			labels: {
				formatter() {
					return `${
						this
							.value
					} min`
				},
			},
			min: 0,
		},

		tooltip: {
			formatter() {
				return `<b>${humanizeDuration(this
					.x)}:</b> ${
					this
						.y
				} min.`
			},
		},

		series: [
			{
				name:
						'Remaining runtime',
				data,
			},
		],
	})
}

function drawUtilityVoltageChart(snapshots) {
	const data = []
	snapshots.forEach((snapshot) => {
		data.push([
			snapshot.duration,
			snapshot.utilityVoltage,
		])
	})

	highcharts.chart({
		...getDefaultChartOptions(),
		chart: {
			renderTo:
					'utility-voltage-chart',
		},

		title: {
			text:
					'Utility voltage',
		},

		yAxis: {
			endOnTick: true,
			labels: {
				formatter() {
					return `${
						this
							.value
					} V`
				},
			},
			min: 0,
		},

		tooltip: {
			formatter() {
				return `<b>${humanizeDuration(this
					.x)}:</b> ${
					this
						.y
				} V`
			},
		},

		series: [
			{
				name:
						'Utility voltage',
				data,
			},
		],
	})
}

function setMonitoringStartTime(startTime) {
	$('#started')
		.append('Monitoring started ')
		.append($(
			'<time />',
			{
				datetime: startTime,
			},
		))

	timeago().render($('time'))
}

function setRevision(revision) {
	$('#revision')
		.append($(
			'<a />',
			{
				href: `https://github.com/jimmysawczuk/power-monitor/commit/${
					revision
						.hex
						.full
				}`,
			},
		).html(`rev. ${
			revision
				.hex
				.short
		}`))
		.append(' &middot; ')
		.append($(
			'<time />',
			{
				datetime:
						revision
							.commit_date
							.iso8601,
			},
		))

	timeago().render($('time'))
}

function update() {
	fetch('/api/snapshots')
		.then(response =>
			response.json())
		.then((response) => {
			if (
				response ===
						null ||
					response.recent ===
						null ||
					response.latest ===
						null
			) {
				return
			}

			const {
				latest,
				recent,
			} = response

			if (
				recent.length ===
					0
			) {
				return
			}

			const now = +moment()
			recent.forEach((
				snapshot,
				i,
			) => {
				recent[
					i
				].duration =
							now -
							+moment(snapshot.timestamp)
			})

			updateModel(latest)

			updateBatteryRemaining(latest)
			updateLoad(latest)
			updateRemainingRuntime(latest)
			updateUtilityVoltage(latest)

			drawBatteryRemainingChart(recent)
			drawLoadChart(recent)
			drawRemainingRuntimeChart(recent)
			drawUtilityVoltageChart(recent)

			$('#last-updated')
				.html('Last updated ')
				.append($(
					'<time />',
					{
						datetime:
									latest.timestamp,
					},
				))

			timeago().render($('time'))
		})
}

function startup(window, document, opts) {
	document.addEventListener(
		'DOMContentLoaded',
		() => {
			window.setInterval(
				update,
				opts.interval,
			)
			update()

			setMonitoringStartTime(opts.startTime)
			setRevision(opts.revision)
		},
	)
}

if (typeof window !== 'undefined') {
	window.startup = startup
}
