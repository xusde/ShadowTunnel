import './App.css';

import {
    ThemeProvider,
} from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import * as React from "react";
import {BarChart, Dashboard, Settings} from "@mui/icons-material";

import {theme} from "./theme";
import {Nav} from "./components/nav/nav";
import {Outlet} from "react-router-dom";
import {useEffect} from "react";
import {GetConfigValue} from "../wailsjs/go/main/App";
import {
    setEncryptionKey,
    setEncryptionMethod,
    setLocalPort,
    setMode,
    setProxyAddress,
    setTransportProtocol
} from "./store/settingsSlice";
import {useDispatch} from "react-redux";


const data = [
    { icon: <Dashboard />, label: 'Dashboard', nav: '/' },
    // { icon: <BarChart />, label: 'Connections', nav: '/connections' },
    { icon: <Settings />, label: 'Settings', nav: '/settings' },
];




function App() {
    const [value, setValue] = React.useState(0);
    const dispatch = useDispatch();

    const handleChange = (event: React.SyntheticEvent, newValue: number) => {
        setValue(newValue);
    };

    useEffect(
        () => {
            GetConfigValue("Mode").then((value) => {
                dispatch(setMode(value));
            });
            GetConfigValue("LocalPort").then((value) => {
                dispatch(setLocalPort(value));
            });
            GetConfigValue("ProxyAddress").then((value) => {
                dispatch(setProxyAddress(value));
            });
            GetConfigValue("EncryptionMethod").then((value) => {
                dispatch(setEncryptionMethod(value));
            });
            GetConfigValue("EncryptionKey").then((value) => {
                dispatch(setEncryptionKey(value));
            });
            GetConfigValue("TransportProtocol").then((value) => {
                dispatch(setTransportProtocol(value));
            });
        },[]
    )


    return (
        <Box id="App">
            <ThemeProvider theme={theme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <Nav data={data}/>
                <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
                    <Outlet />
                </Box>
            </Box>
            </ThemeProvider>
        </Box>
    )
}

export default App
