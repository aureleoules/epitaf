import React from 'react';
import { Redirect, Route, Router as BrowserRouter, Switch } from 'react-router-dom';
import Navbar from './components/Navbar';
import history from './history';
import Home from './routes/Home';
import Groups from './routes/Groups';
import Login from './routes/Login';
import Group from './routes/Group';
import Users from './routes/Users';

import { isLoggedIn } from './utils';

export default function Router() {
	return (
		<BrowserRouter history={history}>
			{isLoggedIn() && (
				<div style={{display: 'flex', minHeight: '100vh'}}>
					<Route>
						<Navbar />
					</Route>

					<Switch>
						<Route exact path="/" component={Home} />
						<Route
							exact
							path="/groups"
							component={Groups}
						/>
						<Route
							exact
							path="/groups/:id"
							component={Group}
						/>

						<Route
							exact
							path="/users"
							component={Users}
						/>
					</Switch>
				</div>
			)}
			{!isLoggedIn() && (
				<Switch>
					<Route exact path="/login" component={Login} />
					<Route
						exact
						component={() => <Redirect to="/login" />}
					/>
				</Switch>
			)}
		</BrowserRouter>
	);
}
