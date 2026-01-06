import { useState, useEffect } from "react";
import { DataGrid } from "@mui/x-data-grid";
import Snackbar from "@mui/material/Snackbar";
import { Button } from "@mui/material";
import { FormControl, InputLabel, Select, MenuItem } from "@mui/material";
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import dayjs from "dayjs";
// import { DataGrid, GridToolbar } from "@mui/x-data-grid";


export default function Report() {
  const baseUrl = sessionStorage.getItem("baseUrl");
  const token = sessionStorage.getItem("token");

  const [value, setValue] = useState("incomes"); // used to select category
  const [from, setFrom] = useState(dayjs().subtract(1, "month")); // set to a month ago   
  const [to, setTo] = useState(dayjs()); // set to current time by default
  

  const [loading, setLoading] = useState(true);
  const [rows, setRows] = useState([]);

  const [snackbarMessage, setSnackbarMessage] = useState("");
  const [showToast, setShowToast] = useState(false);

    const columns = [
    { field: "id", headerName: "ID", width: 70, align: "left", headerAlign: "left" },
    { field: "amount", headerName: "Amount", width: 100, editable: true, type: "number", align: "left", headerAlign: "left" },
    { field: "description", headerName: "Description", width: 300, editable: true, align: "left", headerAlign: "left" },
    { field: "category", headerName: "Category", width: 140, editable: true, align: "left", headerAlign: "left" },
    ];

    const handleReport = async() => {
        // console.log("v = " + value);
        if (value == "incomes") {
            fetchIncomes();
        } else {
            fetchExpenses();
        }
    };

    const handleExport = async () => {

    };

  const fetchIncomes = async () => {
    try {
      setLoading(true);

      const res = await fetch(`${baseUrl}/api/v1/incomes/report?from=${from.unix()}&to=${to.unix()}`, {
        headers: {
          "Authorization": `Bearer ${token}`
        }
      });

      if (!res.ok) throw new Error("Failed to fetch rows");

      const fullData = await res.json();
      const data = fullData || []; // make sure it's an array
      
    //   console.log("data: ", fullData);
  
      setRows(data); // DataGrid expects rows with unique `id`
    } catch (err) {
      setError(err.message || "Error loading rows");
    } finally {
      setLoading(false);
    }
  };

  const fetchExpenses = async () => {
    try {
      setLoading(true);

      const res = await fetch(`${baseUrl}/api/v1/expenses/report?from=${from.unix()}&to=${to.unix()}`, {
        headers: {
          "Authorization": `Bearer ${token}`
        }
      });

      if (!res.ok) throw new Error("Failed to fetch rows");

      const fullData = await res.json();
      const data = fullData || []; // make sure it's an array

      setRows(data); // DataGrid expects rows with unique `id`
    } catch (err) {
      setError(err.message || "Error loading rows");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    handleReport();
  }, []);

  return (
    <div style={{ height: 700, width: "100%" }}>
        
        <div style={{ display:"flex"}}>
            <div style={{ padding:"20px", width:"150px"}}>
                <FormControl fullWidth>
                    <InputLabel>Category</InputLabel>
                    <Select
                        value={value}
                        label="Category"
                        onChange={(e) => setValue(e.target.value)}
                    >
                        <MenuItem value="incomes">Incomes</MenuItem>
                        <MenuItem value="expenses">Expenses</MenuItem>
                    </Select>
                </FormControl>
            </div>

            <div style={{ padding:"20px", width:"300px"}}>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                    <DateTimePicker 
                        label="from"
                        value={from}
                        onChange={(newValue) => setFrom(newValue)}
                    />
                </LocalizationProvider>
            </div>

            <div style={{ padding:"20px", width:"300px"}}>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                    <DateTimePicker 
                        label="to"
                        value={to}
                        onChange={(newValue) => setTo(newValue)}
                        />
                </LocalizationProvider>
            </div>


            <div style={{ padding:"20px", width:"150px"}}>
                <Button variant="contained" onClick={handleReport}>
                    Get report
                </Button>
            </div>

            {/* <div style={{ padding:"20px", width:"150px"}}>
                <Button variant="contained" onClick={handleExport}>
                    Export
                </Button>
            </div> */}


        </div>

        <div style={{ width:"1000px", padding:"20px"}}>
            <DataGrid
                rows={rows}
                columns={columns}
                loading={loading}
                sortModel={[{ field: 'id', sort: 'asc' }]} // initial sort
                pageSizeOptions={[10, 25, 50, 100]}
                // getRowId={(row) => row.ID}
                showToolbar
                // disableColumnFilter
                // disableColumnMenu
                // disableColumnSelector
            />
        </div>


            {/* Snackbar Notification */}
      <Snackbar
        open={showToast}
        autoHideDuration={2000}
        onClose={() => setShowToast(false)}
        message={snackbarMessage}
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
      />

    </div>
  );
}
