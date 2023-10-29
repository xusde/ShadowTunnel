import {createTheme, ThemeOptions} from "@mui/material";

export const themeOptions: ThemeOptions = {
    palette: {
        mode: 'dark',
        primary: {
            main: '#ffb333',
        },
        secondary: {
            main: '#4e342e',
        },
        background: {
            paper: '#3d515f',
            default: '#303030',
        },
        text: {
            primary: '#fff',
        }
    },
    components: {
        MuiListItemButton: {
            defaultProps: {
                disableTouchRipple: true,
            },
        },
        MuiLink: {
            defaultProps: {
                color: 'primary',
            }
        }
    },
    spacing: 8,
};

export const theme = createTheme(themeOptions);