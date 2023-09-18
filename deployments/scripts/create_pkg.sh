echo "Executing create_pkg.sh..."

cd $path_cwd
mkdir $dir_name


cd $source_code_path
# Create and activate virtual environment...
virtualenv -p $runtime env_$function_name
#source $source_code_path/env_$function_name/bin/activate
source env_$function_name/bin/activate

# Installing python dependencies...
#FILE=$source_code_path/requirements.txt
FILE=requirements.txt

if [ -f "$FILE" ]; then
  echo "Installing dependencies..."
  echo "From: requirement.txt file exists..."
  pip install -r "$FILE"

else
  echo "Error: requirement.txt does not exist!"
fi

# Deactivate virtual environment...
deactivate

# Create deployment package...
echo "Creating deployment package..."
cd env_$function_name/lib/$runtime/site-packages/
cp -r . $path_cwd/$dir_name
cd $path_cwd
cp -r $source_code_path/$filename ./$dir_name

# Removing virtual environment folder...
echo "Removing virtual environment folder..."
rm -rf $source_code_path/env_$function_name

echo "Finished script execution!"