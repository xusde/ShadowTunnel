import * as React from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import {MenuItem, Select, SelectChangeEvent, TextField} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../store/store";
import {
    setEncryptionKey,
    setEncryptionMethod,
    setLocalPort,
    setMode,
    setProxyAddress,
    setTransportProtocol
} from "../store/settingsSlice";

export function Settings() {
    const isConnected = useSelector((state: RootState) => state.status.isConnected);
    const mode = useSelector((state:RootState) => state.settings.mode);
    const localPort = useSelector((state:RootState) => state.settings.localPort);
    const proxyAddress = useSelector((state:RootState) => state.settings.proxyAddress);
    const encryptionMethod = useSelector((state:RootState) => state.settings.encryptionMethod);
    const encryptionKey = useSelector((state:RootState) => state.settings.encryptionKey);
    const transportProtocol = useSelector((state:RootState) => state.settings.transportProtocol);

    const dispatch = useDispatch();


    return (
        <List
            // sx={{ width: '100%'}}
            // subheader={<ListSubheader>Settings</ListSubheader>}
        >
            <ListItem>
                <ListItemText id="option-list-label-proxy-mode" primary="Mode" />
                <Select
                    value={mode}
                    onChange={(event) => dispatch(setMode(event.target.value))}
                    size={"small"}
                >
                    <MenuItem value={"direct"}>Direct</MenuItem>
                    <MenuItem value={"proxy"}>Proxy</MenuItem>
                    <MenuItem value={"rules"}>Rules</MenuItem>
                </Select>


            </ListItem>
            <ListItem>
                <ListItemText id="option-list-label-local-port" primary="Local Port" />
                <TextField
                    hiddenLabel
                    value={localPort}
                    onChange={(e) => dispatch(setLocalPort(e.target.value))}
                    disabled={isConnected}
                    variant="outlined"
                    size="small"
                />
            </ListItem>
            <ListItem>
                <ListItemText id="option-list-label-server-addr" primary="Server Address" />
                <TextField
                    hiddenLabel
                    value={proxyAddress}
                    onChange={(e) => dispatch(setProxyAddress(e.target.value))}
                    disabled={isConnected}
                    variant="outlined"
                    size="small"
                />
            </ListItem>
            <ListItem>
                <ListItemText id="option-list-label-encryption-method" primary="Encryption Method" />
                <Select
                    value={encryptionMethod}
                    onChange={(e) => dispatch(setEncryptionMethod(e.target.value))}
                    disabled={isConnected}
                    size={"small"}
                >
                    <MenuItem value={"none"}>None</MenuItem>
                    <MenuItem value={"aes-128-gcm"}>AES-128-GCM</MenuItem>
                </Select>
            </ListItem>
            <ListItem>
                <ListItemText id="option-list-label-encryption-key" primary="Encryption Key" />
                <TextField
                    hiddenLabel
                    value={encryptionKey}
                    onChange={(e) => dispatch(setEncryptionKey(e.target.value))}
                    disabled={isConnected || encryptionMethod === "none"}
                    variant="outlined"
                    size="small"
                />
            </ListItem>
            <ListItem>
                <ListItemText id="option-list-label-transport-protocol" primary="Transport Protocol" />
                <Select
                    value={transportProtocol}
                    onChange={(e) => dispatch(setTransportProtocol(e.target.value))}
                    disabled={isConnected}
                    size={"small"}
                >
                    <MenuItem value={"tcp"}>TCP</MenuItem>
                    <MenuItem value={"quic"}>QUIC</MenuItem>
                </Select>
            </ListItem>

        </List>
    );
}