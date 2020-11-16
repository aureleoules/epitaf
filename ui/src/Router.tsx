import React, { useEffect } from 'react';
import history from './history';

import {Router, Switch, Route, Redirect} from 'react-router-dom';

import Home from './routes/Home';
import Tasks from './routes/Tasks';
import Calendar from './routes/Calendar';
import Callback from './routes/Callback';
import SignIn from './routes/SignIn';
import Sidebar from './components/Sidebar';
import Profile from './routes/Profile';
import { isLoggedIn, logout, parseJwt } from './utils';

export default function(props: any) {

    useEffect(() => {
        const token = localStorage.getItem("jwt");
        if(!token) return;
        try {
            const data = parseJwt(token!);
            if(Date.now() > data.exp * 1000) {
                logout();
            }
        } catch(e) {
            logout();
        }

        const redirect_url = localStorage.getItem("redirect_url");
        if(token && redirect_url) {
            localStorage.setItem("redirect_url", "");
            history.push(redirect_url);
        }
    }, []);
    
    return (
        <Router history={history}>
            <div className="page">
                {isLoggedIn() && <Sidebar/>}
                <Switch>

                    {isLoggedIn() && <>
                        <Route exact path="/" component={Home}/>
                        <Route exact path="/tasks" component={Tasks}/>
                        <Route exact path="/tasks/:id" component={Tasks}/>
                        <Route exact path="/t/:id" component={Tasks}/>
                        <Route exact path="/calendar" component={Calendar}/>
                        <Route exact path="/me" component={Profile}/>

                    </>}
                    {!isLoggedIn() && <>
                        <Route exact path="/login" component={SignIn}/>
                        <Route exact path="/callback" component={Callback}/>
                        <Route render={() => {
                            if(!history.location.pathname.startsWith("/callback")) {
                                const path = history.location.pathname;
                                if(path && path !== "/login" && path !== "/") {
                                    localStorage.setItem("redirect_url", history.location.pathname);
                                }
                                return <Redirect to="/login"/>
                            }
                        }}/>
                    </>}
                </Switch>
            </div>
        </Router>
    )
}