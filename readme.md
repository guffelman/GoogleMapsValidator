# Address Geocoding

This is a command-line tool that takes an Excel file containing addresses and uses the Google Maps Geocoding API to retrieve the latitude and longitude coordinates for each address.

## Installation

1. Clone this repository to your local machine.
2. Install Go on your machine if you haven't already.
3. Run `go build` in the root directory of the project to build the executable file.

## Usage

To use this tool, run the executable file in your terminal and follow the prompts to enter your Google Maps API key and the file path of the Excel file containing the addresses.

The tool expects your unclean addresses to be in column A. It will loop over all addresses at a rate limit of 10/second, check if they are valid using the geocode api, and then split the clean results into columns B-F.


## Contributing

If you find a bug or have a feature request, please open an issue on this repository.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.