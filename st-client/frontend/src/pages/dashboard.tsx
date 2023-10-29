import {Grid, Paper, Stack, styled, Typography, useTheme} from "@mui/material";
import LinkIcon from '@mui/icons-material/Link';
import LinkOffIcon from '@mui/icons-material/LinkOff';
import ArrowCircleUpIcon from '@mui/icons-material/ArrowCircleUp';
import {useSelector} from "react-redux";
import {RootState} from "../store/store";
import {ArrowCircleDown} from "@mui/icons-material";
import {useEffect, useState} from "react";
import {GetSpeed, GetTotalTraffic} from "../../wailsjs/go/main/App";
import ReactEcharts from "echarts-for-react";
import cloneDeep from 'lodash';
import {theme} from "../theme";

const DashboardHeader = styled('h2')(({theme}) => ({
    marginTop: '12px',
    marginLeft: '18px',
    textAlign: 'left',
}));

let downloadSpeedData: number[] = [0,0,0,0,0,0,0,0,0,0];
let uploadSpeedData = [];
let yAxixData = [0, 25, 50, 75, 100, 125, 150, 175, 200];
let xAxisData = ["0s", "1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s"];


const defaultOption = {
    grid: {
        left: "5%",
        right: "5%",
        bottom: "10%",
        top: "10%",
    },
    xAxis: {
        boundaryGap: false,
        splitLine: {
            show: false
        },
        axisLabel: {
            show: true
        },
        axisLine: {
            show: false
        },
        axisTick: {
            show: false
        },
        data: (function (){
            let now:Date = new Date();
            let res = [];
            let len = 10;
            while (len--) {
                res.unshift(now.toLocaleTimeString().replace(/^\D*/,''));
                // @ts-ignore
                now = new Date(now - 2000);
            }
            return res;
        })()
    },
    yAxis: {
        boundaryGap: false,
        splitLine: {
            show: false
        },
        axisLabel: {
            color: theme.palette.text.primary
        },
        data: yAxixData
    },
    series: [
        {
            data: downloadSpeedData,
            type: 'line',
            smooth: true,
            areaStyle: {},
        }
    ]
};

export function Dashboard() {
    const theme = useTheme();
    const isConnected = useSelector((state: RootState) => state.status.isConnected);
    const [uploadSpeed, setUploadSpeed] = useState("0.00");
    const [downloadSpeed, setDownloadSpeed] = useState("0.00");
    const [totalTraffic, setTotalTraffic] = useState("0.00");
    const [option, setOption] = useState(defaultOption);

    function updateSpeedDiagram(uploadSpeed: number, downloadSpeed: number) {
        let newOption = defaultOption;

        downloadSpeedData.push(downloadSpeed);
        // uploadSpeedData.push(uploadSpeed);
        downloadSpeedData.shift();
        // uploadSpeedData.shift();
        console.log(downloadSpeedData);

        xAxisData.shift();
        xAxisData.push(new Date().toLocaleTimeString().replace(/^\D*/,''));

        console.log(xAxisData);

        newOption.series[0].data = downloadSpeedData;
        // newOption.series[1].data = uploadSpeedData;
        newOption.xAxis.data = xAxisData;

        console.log(newOption);
        //
        // // @ts-ignore
        // newOption.series[0].data.shift();
        // // @ts-ignore
        // newOption.series[0].data.push(downloadSpeed);
        // // @ts-ignore
        // newOption.xAxis.data.shift();
        // // @ts-ignore
        // newOption.xAxis.data.push(new Date().toLocaleTimeString().replace(/^\D*/,''));
        // // @ts-ignore
        setOption(newOption);
    }

    useEffect(() => {
        // update speed every 1 second
        const interval = setInterval(() => {
            console.log('update speed');
            GetSpeed("upload").then((speed) => {
                // console.log("upload " + speed);
                setUploadSpeed(speed);
            })
            GetSpeed("download").then((speed) => {
                // console.log("download: " + speed);
                setDownloadSpeed(speed);
                console.log("download: " + parseFloat(speed));
                updateSpeedDiagram(0, parseFloat(speed));
            })
            GetTotalTraffic().then((traffic) => {
                // console.log("total traffic: " + traffic);
                setTotalTraffic(traffic);
            })
        }, 3000);
        return () => clearInterval(interval);
    }, []);

    return (
        <Grid container spacing={2}>
            <Grid item xs={4}>
                <Paper sx={{p: 2, display: 'flex', flexDirection: 'column', height: 240, textAlign: 'center', backgroundColor: isConnected ? theme.palette.success.dark: theme.palette.background.paper}}>
                    <DashboardHeader>Status</DashboardHeader>
                    {isConnected ? <LinkIcon sx={{fontSize: '48px', textAlign: 'center', width: '100%', height: '50%'}}/> : <LinkOffIcon sx={{fontSize: '48px', textAlign: 'center', width: '100%', height: '50%'}}/>}
                </Paper>
            </Grid>
            <Grid item xs={4}>
                <Paper sx={{p: 2, display: 'flex', flexDirection: 'column', height: 240}}>
                    <DashboardHeader>Speed</DashboardHeader>
                    <Stack direction={"row"} alignItems={"center"} justifyContent={"center"} spacing={1} sx={{paddingBottom: '8px', paddingTop: '14px'}}>
                        <ArrowCircleUpIcon sx={{fontSize: '36px', textAlign: 'center'}}/>
                        <Typography variant="h1" component="div" sx={{textAlign: 'center', fontSize: '28px', fontWeight: 'bold'}}>
                            {isConnected ? uploadSpeed : "-.--"} KB/s
                        </Typography>
                    </Stack>
                    <Stack direction={"row"} alignItems={"center"} justifyContent={"center"} spacing={1}>
                        <ArrowCircleDown sx={{fontSize: '36px', textAlign: 'center'}}/>
                        <Typography variant="h1" component="div" sx={{textAlign: 'center', fontSize: '28px', fontWeight: 'bold'}}>
                            {isConnected ? downloadSpeed : "-.--"} KB/s
                        </Typography>
                    </Stack>
                </Paper>
            </Grid>
            <Grid item xs={4}>
                <Paper sx={{p: 2, display: 'flex', flexDirection: 'column', height: 240}}>
                    <DashboardHeader>Traffic</DashboardHeader>
                    <Typography variant="h1" component="div" sx={{fontSize: '32px', textAlign: 'center', fontWeight: 'bold', paddingTop: '28px'}}>
                        {totalTraffic} MB
                    </Typography>
                </Paper>
            </Grid>
            <Grid item xs={12}>
                <Paper sx={{p: 2, display: 'flex', flexDirection: 'column', height: 295}}>
                    {/*<ReactEcharts option={option}/>*/}
                </Paper>
            </Grid>
        </Grid>
    )
}