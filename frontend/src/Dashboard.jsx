import PieChartExpenses from "./DonutchartExpenses";
import PieChartIncomes from "./DonutchartIncomes";
import PieChartSummary from "./DonutchartSummary";
import Typography from "@mui/material/Typography";

export default function Dashboard() {

  return (
        <div style={{ padding: "20px", display: "flex"}}>
            
            <div style={{padding:"50px"}}>
                <PieChartExpenses />
                <Typography
                variant="body2"
                align="center"
                sx={{ mt: 1, color: "text.secondary" }}
                >
                Expenses Summary
              </Typography>
            </div>
            <div style={{padding:"50px"}}>
              <PieChartIncomes />
                <Typography
                  variant="body2"
                  align="center"
                  sx={{ mt: 1, color: "text.secondary" }}
                >
                  Incomes Summary
                </Typography>
            </div>

            <div style={{padding:"50px"}}>
              <PieChartSummary />
                <Typography
                  variant="body2"
                  align="center"
                  sx={{ mt: 1, color: "text.secondary" }}
                >
                  Summary
                </Typography>
            </div>
            
        </div>
  );
}