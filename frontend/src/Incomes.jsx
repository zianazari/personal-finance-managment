import { useState, useEffect } from "react";
import { DataGrid } from "@mui/x-data-grid";
import Snackbar from "@mui/material/Snackbar";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from "@mui/material";

export default function Incomes() {
  const baseUrl = sessionStorage.getItem("baseUrl");
  const token = sessionStorage.getItem("token");

  const [loading, setLoading] = useState(true);
  const [rows, setRows] = useState([]);

  const [open, setOpen] = useState(false);
  const [formData, setFormData] = useState({
    amount: "",
    description: "",
    category: "",
  });

  const [snackbarMessage, setSnackbarMessage] = useState("");
  const [showToast, setShowToast] = useState(false);

  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);
  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const columns = [
    { field: "displayId", headerName: "ID", width: 70, align: "left", headerAlign: "left" },
    { field: "amount", headerName: "Amount", width: 100, editable: true, type: "number", align: "left", headerAlign: "left" },
    { field: "description", headerName: "Description", width: 300, editable: true, align: "left", headerAlign: "left" },
    { field: "category", headerName: "Category", width: 140, editable: true, align: "left", headerAlign: "left" },
  {
    field: "actions",
    headerName: "Actions",
    width: 250,
    renderCell: (params) => (
      <div style={{ display: "flex", gap: "10px" }}>
        <Button
          variant="contained"
          color="primary"
          onClick={() => handleUpdate(params.row)}
        >
          Update
        </Button>

        <Button
          variant="contained"
          color="error"
          onClick={() => handleDelete(params.row.id)}
        >
          Delete
        </Button>
      </div>
    ),
    },
  ];

  const fetchIncomes = async () => {
    try {
      setLoading(true);

      const res = await fetch(`${baseUrl}/api/v1/incomes/list`, {
        headers: {
          "Authorization": `Bearer ${token}`
        }
      });

      if (!res.ok) throw new Error("Failed to fetch incomes");

      const fullData = await res.json();
      const data = fullData.ok || []; // make sure it's an array

      // setRows(data); // DataGrid expects rows with unique `id`
      const processed = data.map((row, index) => ({
      ...row,
      displayId: index + 1,
    }));

    setRows(processed);

    } catch (err) {
      setError(err.message || "Error loading incomes");
    } finally {
      setLoading(false);
    }
  };

  const handleInsert = async() => {
    const newRecord = {
      amount: Number(formData.amount),
      description: formData.description,
      category: formData.category,
    };
      try {
      // const payload = { ...row, amount: parseInt(row.amount, 10) };

      const res = await fetch(`${baseUrl}/api/v1/incomes/add`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(newRecord), // Send entire row info
      });

      const data = await res.json();
      
      if (data.error) {
        setSnackbarMessage(data.error);
        setShowToast(true);
      }
      
      setSnackbarMessage(data.ok || "Row added successfully!");
      setShowToast(true);

      // Refresh rows
      fetchIncomes();

      // Close popup and reset
      setFormData({ amount: "", description: "", category: "" });
      handleClose();
      
    } catch (err) {
      setSnackbarMessage(err.messaage);
      setShowToast(true);
    }
  };

  const handleDelete = async(id) => {
      try {
   
      const res = await fetch(`${baseUrl}/api/v1/incomes/delete?id=${id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        // body: JSON.stringify(payload), // Send entire row info
      });

      const data = await res.json();
      
      if (data.error) {
        setSnackbarMessage(data.error);
        setShowSnackbar(true);
      }
      
      setSnackbarMessage(data.ok || "Row deleted successfully!");
      setShowToast(true);

      // Refresh rows
      fetchIncomes();
      
    } catch (err) {
      setSnackbarMessage(err.messaage);
      setShowSnackbar(true);
    }
  };

  const handleUpdate = async(row) => {
      try {
      const payload = { ...row, amount: parseInt(row.amount, 10) };

      const res = await fetch(`${baseUrl}/api/v1/incomes/update`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload), // Send entire row info
      });

      const data = await res.json();
      
      if (data.error) {
        setSnackbarMessage(data.error);
        setShowSnackbar(true);
      }
      
      setSnackbarMessage(data.ok || "Row updated successfully!");
      setShowToast(true);

      // Refresh rows
      fetchIncomes();
      
    } catch (err) {
      setSnackbarMessage(err.messaage);
      setShowSnackbar(true);
    }
  };

  useEffect(() => {
    fetchIncomes();
  }, []);

  return (
    <div style={{ height: 600, width: "100%" }}>

      <div style={{ padding:"20px", width:"200px"}}>
        <Button variant="contained" onClick={handleOpen}>
          Add New Income
        </Button>
      </div>

      {/* POPUP FORM */}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Add a new income record</DialogTitle>
        <DialogContent sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
          <TextField
            label="Amount"
            name="amount"
            type="number"
            value={formData.amount}
            onChange={handleChange}
            fullWidth
          />
          <TextField
            label="Description"
            name="description"
            value={formData.description}
            onChange={handleChange}
            fullWidth
          />
          <TextField
            label="Category"
            name="category"
            value={formData.category}
            onChange={handleChange}
            fullWidth
          />
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button variant="contained" onClick={handleInsert}>
            Send
          </Button>
        </DialogActions>
      </Dialog>


    <div style={{ width:"1000px", padding:"20px"}}>
          <DataGrid
            rows={rows}
            columns={columns}
            loading={loading}
            sortModel={[{ field: 'id', sort: 'asc' }]} // initial sort
            pageSizeOptions={[10, 25, 50]}
            showToolbar
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
