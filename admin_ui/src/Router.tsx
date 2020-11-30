import React from 'react';
import history from './history';

import {Router as BrowserRouter, Switch, Route} from 'react-router-dom';

import Home from './routes/Home';

export default function Router(props: any) {
    return (
        <BrowserRouter history={history}>
            <Switch>
                <Route exact path="/" component={Home}/>
            </Switch>
        </BrowserRouter>
    )
}