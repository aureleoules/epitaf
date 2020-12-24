import React from 'react';
import history from './history';

import {Router as BrowserRouter, Switch, Route, Redirect} from 'react-router-dom';
import {isLoggedIn} from './utils';

import Home from './routes/Home';
import Login from './routes/Login';
import Tasks from './routes/Tasks';
import Task from './routes/Task';
import Groups from './routes/Groups';
import Group from './routes/Group';
import Users from './routes/Users';
import User from './routes/User';
import EditUser from './routes/EditUser';
import Roles from './routes/Roles';
import Role from './routes/Role';
import Admins from './routes/Admins';
import Admin from './routes/Admin';
import Profile from './routes/Profile';
import Settings from './routes/Settings';
import Drawer from './components/Drawer';

export default function Router(props: any) {
    return (
        <BrowserRouter history={history}>
            {isLoggedIn() && <Drawer>
            <Switch>
                    <Route exact path="/" component={Home}/>

                    <Route exact path="/tasks" component={Tasks}/>
                    <Route exact path="/tasks/:id" component={Task}/>

                    <Route exact path="/groups" component={Groups}/>
                    <Route exact path="/groups/:id" component={Group}/>
                    
                    <Route exact path="/users" component={Users} />
                    <Route exact path="/users/new" component={(props: any) => <EditUser new {...props}/>}/>
                    <Route exact path="/users/:id" component={User}/>
                    <Route exact path="/users/:id/edit" component={EditUser}/>

                    <Route exact path="/roles" component={Roles}/>
                    <Route exact path="/roles/:id" component={Role}/>

                    <Route exact path="/admins" component={Admins}/>
                    <Route exact path="/admins/:id" component={Admin}/>

                    <Route exact path="/profile" component={Profile}/>
                    <Route exact path="/settings" component={Settings}/>

                </Switch>
            </Drawer>}

            {!isLoggedIn() && <Switch>
                <Route exact path="/login" component={Login}/>
                <Route exact path="/" component={() => <Redirect to="/login"/>}/>
            </Switch>}
        </BrowserRouter>
    )
}