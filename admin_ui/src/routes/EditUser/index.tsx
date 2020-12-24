import React, { useState } from 'react';

import ButtonGroup from '@material-ui/core/ButtonGroup';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import { useTranslation } from 'react-i18next';
import { createStyles, makeStyles, Theme } from '@material-ui/core';
import GroupTree from '../../views/GroupTree';


const useStyles = makeStyles((theme: Theme) => createStyles({
	root: {
		display: 'flex',
		flexWrap: 'wrap',
		flexDirection: "column",
	},
	textField: {
		width: 350,
		marginBottom: 15
	},
	title: {
		marginBottom: 5
	}
}));

export default function EditUser(props: any) {
    
    const [name, setName] = useState<string>("");
    const [email, setEmail] = useState<string>("");
    const [login, setLogin] = useState<string>("");
    const [password, setPassword] = useState<string>("");

    const classes = useStyles();
        
    const {t} = useTranslation();

    return (
        <>
            <h1>Edit user</h1>

            <form className={classes.root}>
                <Typography className={classes.title} variant="h5">Informations</Typography>
                <div>
                    <TextField
                        className={classes.textField}
                        required
                        label={t('Fullname')}
                        onChange={e => setName(e.target.value)}
                        value={name}
                    />
                </div>
                <div>
                    <TextField
                        className={classes.textField}
                        required
                        label={t('Email')}
                        onChange={e => setEmail(e.target.value)}
                        value={email}
                    />
                </div>
                <div>
                    <TextField
                        className={classes.textField}
                        required
                        label={t('Email')}
                        onChange={e => setEmail(e.target.value)}
                        value={email}
                        type="email"
                    />
                </div>
                <div>
                    <TextField
                        className={classes.textField}
                        required
                        label={t('Login')}
                        onChange={e => setLogin(e.target.value)}
                        value={login}
                    />
                </div>
                <div>
                    <TextField
                        className={classes.textField}
                        required
                        label={t('Password')}
                        onChange={e => setPassword(e.target.value)}
                        value={password}
                        type="password"
                    />
                </div>
                <div>
                    <Button color="primary" variant={"contained"}>{t('Save')}</Button>
                </div>

            </form>
            <GroupTree/>
        </>
    )
}