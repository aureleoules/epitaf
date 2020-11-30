import AppBar from "@material-ui/core/AppBar";
import CssBaseline from "@material-ui/core/CssBaseline";
import DrawerMenu from "@material-ui/core/Drawer";
import IconButton from "@material-ui/core/IconButton";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import AccountBox from "@material-ui/icons/AccountBox";
import Dashboard from "@material-ui/icons/Dashboard";
import MenuIcon from "@material-ui/icons/Menu";
import Settings from "@material-ui/icons/Settings";
import React, { useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import Link from "@material-ui/core/Link";
import AccountCircle from "@material-ui/icons/AccountCircle";
import Menu from "@material-ui/core/Menu";
import { useTranslation } from 'react-i18next';
import { ClickAwayListener } from "@material-ui/core";

const drawerWidth = 240;

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            display: "flex",
        },
        appBar: {
            zIndex: theme.zIndex.drawer + 1,
            transition: theme.transitions.create(["width", "margin"], {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.leavingScreen,
            }),
        },
        appBarShift: {
            marginLeft: drawerWidth,
            width: `calc(100% - ${drawerWidth}px)`,
            transition: theme.transitions.create(["width", "margin"], {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
            }),
        },
        menuButton: {
            marginRight: 36,
        },
        hide: {
            display: "none",
        },
        drawer: {
            width: drawerWidth,
            flexShrink: 0,
        },
        drawerOpen: {
            width: drawerWidth,
            transition: theme.transitions.create("width", {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
            }),
        },
        drawerClose: {
            transition: theme.transitions.create("width", {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.leavingScreen,
            }),
            overflowX: "hidden",
            width: 80,
            [theme.breakpoints.up("sm")]: {
                width: 80,
            },
        },
        drawerPaper: {
            borderRight: "none",
            overflow: "hidden",
        },
        drawerContainer: {
            overflow: "auto",
        },
        toolbar: {
            display: "flex",
            alignItems: "center",
            justifyContent: "flex-end",
            padding: theme.spacing(0, 1),
            // necessary for content to be below app bar
            ...theme.mixins.toolbar,
        },
        listItem: {
            paddingLeft: 28,
        },
        content: {
            flexGrow: 1,
            padding: theme.spacing(3),
        },
        bottomItems: {
            position: "absolute",
            bottom: 0,
            width: "100%",
        },
        itemLink: {
            color: "inherit",
            textDecoration: "none",
            "&:hover": {
                textDecoration: "none",
            },
        },
        userMenu: {
            position: "absolute",
            right: 15,
        },
    })
);

type Props = {
    title?: string
    children?: any
    menu?: any
    menuIcon?: any
}

export default function Drawer(props: Props) {
    const classes = useStyles();
    const [open, setOpen] = useState<boolean>(true);

    const { t } = useTranslation();

    const handleDrawerToggle = () => {
        setOpen(!open);
    };

    const [menu, setMenu] = useState<boolean>(false);
    function openMenu() {
        setMenu(true);
    }
    function closeMenu() {
        setMenu(false);
    }

    const routes = [
        {
            name: "Home",
            icon: Dashboard,
            path: "/",
        }
    ];

    const bottomRoutes = [
        {
            name: "Users",
            icon: AccountBox,
        },
        {
            name: "Settings",
            icon: Settings,
        },
    ];

    return (
        <div className={classes.root}>
            <CssBaseline />
            <AppBar position="fixed" className={clsx(classes.appBar)}>
                <Toolbar>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        onClick={handleDrawerToggle}
                        edge="start"
                        className={clsx(classes.menuButton)}
                    >
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" noWrap>
                        {props.title || t('Dashboard')}
                    </Typography>
                    <div className={classes.userMenu}>
                        <ClickAwayListener onClickAway={closeMenu}>
                            <IconButton
                                onClick={openMenu}
                                aria-label="account of current user"
                                aria-controls="menu-appbar"
                                aria-haspopup="true"
                                color="inherit"
                            >
                                {props.menuIcon ? <props.menuIcon /> : <AccountCircle />}
                            </IconButton>
                        </ClickAwayListener>
                        <Menu
                            id="menu-appbar"
                            anchorOrigin={{
                                vertical: "top",
                                horizontal: "right",
                            }}
                            keepMounted
                            transformOrigin={{
                                vertical: "top",
                                horizontal: "right",
                            }}
                            open={menu}
                        >
                            {props.menu}
                        </Menu>
                    </div>
                </Toolbar>
            </AppBar>
            <DrawerMenu
                className={classes.drawer}
                variant="permanent"
                classes={{
                    paper: clsx(classes.drawerPaper, {
                        [classes.drawerClose]: !open,
                        [classes.drawerOpen]: open,
                    }),
                }}
            >
                <div className={classes.toolbar} />
                <List>
                    {routes.map((r) => (
                        <Link
                            className={classes.itemLink}
                            component={RouterLink}
                            to={r.path}
                            key={r.name}
                        >
                            <ListItem className={classes.listItem} button >
                                <ListItemIcon>
                                    <r.icon />
                                </ListItemIcon>
                                <ListItemText primary={r.name} />
                            </ListItem>
                        </Link>
                    ))}
                </List>
                <List className={classes.bottomItems}>
                    {bottomRoutes.map((r) => (
                        <ListItem className={classes.listItem} button key={r.name}>
                            <ListItemIcon>
                                <r.icon />
                            </ListItemIcon>
                            <ListItemText primary={r.name} />
                        </ListItem>
                    ))}
                </List>
            </DrawerMenu>
            <main className={classes.content}>
                <div className={classes.toolbar} />
                {props.children}
            </main>
        </div>
    );
}