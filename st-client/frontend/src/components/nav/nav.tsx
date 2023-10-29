import Button from "@mui/material/Button";
import {ListItem, ListItemButton, styled} from "@mui/material";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import * as React from "react";
import MuiDrawer from "@mui/material/Drawer";
import {Link} from "react-router-dom";
import {themeOptions} from "../../theme";
import {useDispatch, useSelector} from "react-redux";
import {toggleConnection} from "../../store/statusSlice";
import {RootState} from "../../store/store";

import {Connect, Disconnect} from "../../../wailsjs/go/main/App";

const drawerWidth = 180;

const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    () => ({
        width: drawerWidth,
        // flexShrink: 0,
        // whiteSpace: 'nowrap',
        // boxSizing: 'border-box',
    }),
);

export interface DrawerProps {
    data: { icon: React.ReactNode; label: string; nav: string }[];
}

export function Nav(props: DrawerProps) {
    const isConnected = useSelector((state: RootState) => state.status.isConnected);
    const localPort = useSelector((state:RootState) => state.settings.localPort);
    const proxyAddress = useSelector((state:RootState) => state.settings.proxyAddress);
    const dispatch = useDispatch();

    const handleToggleConnection = () => {
        dispatch(toggleConnection());
        if (isConnected) {
            Disconnect();
        } else {
            Connect(`localhost:${localPort}`, proxyAddress);
        }
    }

    return (
        <Drawer variant="permanent" sx={{
            display: 'block',
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth, paddingTop: '36px', border: 'none' },
        }}>
            <Button onClick={handleToggleConnection} variant="contained" color={'primary'} sx={{height: '36px', paddingX: '0px', marginX: '16px', marginBottom: '24px'}}>
                {isConnected ? 'Disconnect' : 'Connect'}
            </Button>
            {props.data.map((item) => (
                <ListItem sx={{ color: themeOptions.palette?.text?.primary }} key={item.label} disablePadding component={Link} to={item.nav}>
                    <ListItemButton
                        sx={{ py: '4px', minHeight: 32}}
                    >
                        <ListItemIcon sx={{ color: 'inherit', paddingX: '8px', marginRight: '-8px' }}>
                            {item.icon}
                        </ListItemIcon>
                        <ListItemText
                            primary={item.label}
                            primaryTypographyProps={{ fontSize: 16, fontWeight: 'medium' }}
                        />
                    </ListItemButton>
                </ListItem>
            ))}
        </Drawer>
    )
}