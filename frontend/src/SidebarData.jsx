import HomeIcon from '@mui/icons-material/Home';
import EuroIcon from '@mui/icons-material/Euro';
import ShopIcon from '@mui/icons-material/Shop';
import LogoutIcon from '@mui/icons-material/Logout';
import GroupIcon from '@mui/icons-material/Group';
import AssessmentIcon from '@mui/icons-material/Assessment';

export const SidebarData = [
    {
        title: "home",
        icon: <HomeIcon />,
        link: "/"
    },
    {
        title: "Incomes",
        icon: <EuroIcon />,
        link: "/incomes"
    },
    {
        title: "Expenses",
        icon: <ShopIcon />,
        link: "/expenses"
    },
    {
        title: "Report",
        icon: <AssessmentIcon />,
        link: "/report"
    },
    {
        title: "Users",
        icon: <GroupIcon />,
        link: "/users"
    },
    {
        title: "Logout",
        icon: <LogoutIcon />,
        link: "/logout"
    },
]