import React, { Component } from 'react'
import axios from 'axios';

class ApplicationDashboard extends Component {
	
	state = {
		avgRespTime : null,
		totalRequests : null,
		childStatus : []
	}

	componentDidMount() {
		axios.get('http://localhost:4000/api/analytics')
			.then(res => {

				const responseData = res["data"]
				console.log(responseData)
				this.setState({
					avgRespTime : responseData["average_response_time"],
					childStatus : responseData["stats_per_route"],
					totalRequests : responseData["total_requests"]
				})
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
