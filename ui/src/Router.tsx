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
import { getTheme, isLoggedIn, logout, parseJwt } from './utils';
import { Console } from 'console';

export default function(props: any) {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [_, y] = React.useState(false);
    const [ok, setOk] = React.useState<boolean>(false);

    useEffect(() => {
        setTimeout(() => {
            setOk(true);
        }, 10);
        
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

        window.addEventListener("render", () => {
            y(x => !x);
        });

    }, []);
    
    return (
        <Router history={history}>
            <div className={[!ok ? "preload": "", "router", "page", "theme", "theme--" + getTheme()].join(" ")}>
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