import * as React from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';
import ImageIcon from '@mui/icons-material/Image';
import {EventsOn} from "../../wailsjs/runtime";
import {Paper} from "@mui/material";

export function Connections() {
    const initalConnections:string[] = [];
    const [connections, setConnections] = React.useState(initalConnections);
    EventsOn("updateConnList", (data:string) => {
        console.log(data);
        // setConnections(data);
        setConnections((connections:string[]) => [...connections, data]);
    });
    return (
        <Paper sx={{paddingY: '12px'}}>
            <List sx={{ width: '100%' }}>
                {/*{connections.map((connection) => {*/}
                {/*    return (*/}
                {/*        <ListItem key={connection.id}>*/}
                {/*            <ListItemAvatar>*/}
                {/*                <Avatar>*/}
                {/*                    <ImageIcon />*/}
                {/*                </Avatar>*/}
                {/*            </ListItemAvatar>*/}
                {/*            <ListItemText primary={connection.name} secondary={connection.id} />*/}
                {/*        </ListItem>*/}
                {/*    )*/}
                {/*})}*/}
                {
                    connections.length !==0 ? connections.map((connection) => {
                    return (
                        <ListItem key={connection}>
                            <ListItemAvatar>
                                <Avatar>
                                    <ImageIcon />
                                </Avatar>
                            </ListItemAvatar>
                            <ListItemText primary={connection}/>
                        </ListItem>
                    )
                })
                    : <h3>No connections</h3>
                }



                {/*<ListItem>*/}
                {/*    <ListItemAvatar>*/}
                {/*        <Avatar>*/}
                {/*            <ImageIcon />*/}
                {/*        </Avatar>*/}
                {/*    </ListItemAvatar>*/}
                {/*    <ListItemText primary="Photos" secondary="Jan 9, 2014" />*/}
                {/*</ListItem>*/}
            </List>
        </Paper>
    );
}

