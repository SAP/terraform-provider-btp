# terraform import btp_directory.<resource_name> <directory_id>

terraform import btp_directory.parent dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0

#terraform import using id attribute in import block

import {
  to = btp_directory.<resource_name>
  id = "<directory_id>"
}
