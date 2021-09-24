/*
Copyright © 2021 SignorMercurio

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package msc

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/spf13/cobra"
)

// NewJpgCmd represents the jpg command
func NewJpgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jpg",
		Short: "Read EXIF info in JPG files",
		Long: `Read EXIF info in JPG files
Example:
	att msc -i in.jpg -o out.txt jpg
	att msc -i in.tiff jpg`,
		RunE: func(cmd *cobra.Command, args []string) error {
			x, err := exif.Decode(input)
			if err != nil {
				if exif.IsCriticalError(err) {
					return errors.Wrap(err, "decode EXIF info")
				} else {
					Warn("Error decoding some of EXIF info， trying to proceed...")
				}
			}
			Echo(x.String())
			lat, long, err := x.LatLong()
			if err != nil {
				Warn("Error extracting geo info")
			} else {
				Echo(fmt.Sprintf("Location: (%f, %f)\n", lat, long))
			}

			return nil
		},
	}

	return cmd
}

func init() {
	mscCmd.AddCommand(NewJpgCmd())
}
