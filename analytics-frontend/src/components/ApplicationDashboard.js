import React, { Component } from 'react'
import axios from 'axios';
import Pusher from 'pusher-js'

class ApplicationDashboard extends Component {
	
	state = {
		avgRespTime : null,
		totalRequests : null,
		childStatus : []
	}

	componentDidMount() {
		this.setUpPusher()

		axios.get('http://localhost:4000/api/analytics')
			.then(res => {
				const responseData = res["data"]
				this.setMetrics(responseData)
			})
	}

	// Updating state variables with the response data
	setMetrics = (responseData) => {
		this.setState({
			avgRespTime : responseData["average_response_time"],
			childStatus : responseData["stats_per_route"],
			totalRequests : responseData["total_requests"]
		})
	}

	// Setting up pusher subscribers for updating the UI based on events from publisher
	setUpPusher = () => {
		// environment variables
		const app_key = process.env.REACT_APP_PUSHER_APP_KEY
		const cluster_key = process.env.REACT_APP_PUSHER_APP_CLUSTER
		// setting up pusher object
		var pusher = new Pusher(app_key, {
			cluster: cluster_key
		})
		// subscribing to the channel grpc-monitoring
		var channel = pusher.subscribe("grpc-monitoring")
		channel.bind('data', data => {
			this.setMetrics(data)
		})
	}
	
	render() {
		return (
			<div>
				<h1>Total Requests : {this.state.totalRequests}</h1>
				<h2>Average response time : {this.state.avgRespTime}</h2>
				<ul>
					{this.state.childStatus.map( stats => 
						<li key={stats.id.url}><b>URL:</b> {stats.id.url} <b>Number of Requests:</b> {stats.number_of_requests}</li>
					)}
				</ul>
			</div>
		)
	}
}

export default ApplicationDashboard
