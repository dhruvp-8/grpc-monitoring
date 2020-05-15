import React, { Component } from 'react'
import axios from 'axios';

class ApplicationDashboard extends Component {
	
	componentDidMount() {
		axios.get('http://localhost:4000/api/analytics')
			.then(res => {
				console.log(res["data"])
			})
	}
	
	render() {
		return (
			<div>
				<h1>Count : 1</h1>
			</div>
		)
	}
}

export default ApplicationDashboard
