import React from 'react';
import history from './history';

import {Router as BrowserRouter, Switch, Route, Redirect} from 'react-router-dom';
import {isLoggedIn} from './utils';

import Home from './routes/Home';
import Login from './routes/Login';

export default function Router(props: any) {
    return (
        <BrowserRouter history={history}>
            <Switch>
                {isLoggedIn() && <>
                    <Route exact path="/" component={Home}/>
                </>}

                {!isLoggedIn() && <>
                    <Route exact path="/login" component={Login}/>
                    <Route exact path="/" component={() => <Redirect to="/login"/>}/>
                </>}
            </Switch>
        </BrowserRouter>
    )
}