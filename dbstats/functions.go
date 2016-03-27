package dbstats

import (
	"database/sql"

	null "gopkg.in/guregu/null.v2"
)

type Function struct {
	SchemaName      string      `json:"schema_name"`
	FunctionName    string      `json:"function_name"`
	Language        string      `json:"language"`
	Source          string      `json:"source"`
	SourceBin       null.String `json:"source_bin"`
	Config          null.String `json:"config"`
	Arguments       null.String `json:"arguments"`
	Result          null.String `json:"result"`
	Aggregate       bool        `json:"aggregate"`
	Window          bool        `json:"window"`
	SecurityDefiner bool        `json:"security_definer"`
	Leakproof       bool        `json:"leakproof"`
	Strict          bool        `json:"strict"`
	ReturnsSet      bool        `json:"returns_set"`
	Volatile        null.String `json:"volatile"`
	Calls           null.Int    `json:"calls"`
	TotalTime       null.Float  `json:"total_time"`
	SelfTime        null.Float  `json:"self_time"`
}

const functionsSQL string = `
SELECT pn.nspname AS schema_name,
			 pp.proname AS function_name,
			 pl.lanname AS language,
			 pp.prosrc AS source,
			 pp.probin AS source_bin,
			 pp.proconfig AS config,
			 pg_get_function_arguments(pp.oid) AS arguments,
			 pg_get_function_result(pp.oid) AS result,
			 pp.proisagg AS aggregate,
			 pp.proiswindow AS window,
			 pp.prosecdef AS security_definer,
			 pp.proleakproof AS leakproof,
			 pp.proisstrict AS strict,
			 pp.proretset AS returns_set,
			 pp.provolatile AS volatile,
			 ps.calls,
			 ps.total_time,
			 ps.self_time
	FROM pg_proc pp
 INNER JOIN pg_namespace pn ON (pp.pronamespace = pn.oid)
 INNER JOIN pg_language pl ON (pp.prolang = pl.oid)
	LEFT JOIN pg_stat_user_functions ps ON (ps.funcid = pp.oid)
 WHERE pl.lanname != 'internal'
			 AND pn.nspname NOT IN ('pg_catalog', 'information_schema')
			 AND pp.proname NOT IN ('pg_stat_statements', 'pg_stat_statements_reset')`

func GetFunctions(db *sql.DB, postgresVersionNum int) ([]Function, error) {
	stmt, err := db.Prepare(QueryMarkerSQL + functionsSQL)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var functions []Function

	for rows.Next() {
		var row Function

		err := rows.Scan(&row.SchemaName, &row.FunctionName, &row.Language, &row.Source,
			&row.SourceBin, &row.Config, &row.Arguments, &row.Result, &row.Aggregate,
			&row.Window, &row.SecurityDefiner, &row.Leakproof, &row.Strict, &row.ReturnsSet,
			&row.Volatile, &row.Calls, &row.TotalTime, &row.SelfTime)
		if err != nil {
			return nil, err
		}

		functions = append(functions, row)
	}

	return functions, nil
}