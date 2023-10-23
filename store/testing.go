package tuples

import "context"


func Cleanup() {
	ctx := context.Background()
	pg.Exec(ctx, `delete from tuples`)
}
